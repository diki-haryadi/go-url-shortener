package shortener

import (
	"context"
	"errors"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/go-redis/redis/v8"
	"go-url-shortener/helpers"
	"go-url-shortener/pkg/cache"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

// Service is a simple CRUD interface for user profiles.
type Service interface {
	PostShorten(ctx context.Context, p Shorten) (shorten ShortenResp, err error)
	Resolve(ctx context.Context, url string) (string, error)
}

// Shorten represents a single user profile.
// ID should be globally unique.
type Shorten struct {
	URL         string        `json:"url"`
	CustomShort string        `json:"short"`
	Expiry      time.Duration `json:"expiry"`
}

type ShortenResp struct {
	URL             string        `json:"url"`
	CustomShort     string        `json:"short"`
	Expiry          time.Duration `json:"expiry"`
	XRateRemaining  int           `json:"rate_limit"`
	XRateLimitReset time.Duration `json:"rate_limit_reset"`
}

// Address is a field of a user profile.
// ID should be unique within the profile (at a minimum).
type Address struct {
	ID       string `json:"id"`
	Location string `json:"location,omitempty"`
}

var (
	ErrInconsistentIDs = errors.New("inconsistent IDs")
	ErrAlreadyExists   = errors.New("already exists")
	ErrNotFound        = errors.New("not found")
)

type inmemService struct {
	mtx sync.RWMutex
	m   map[string]Shorten
}

func NewInmemService() Service {
	return &inmemService{
		m: map[string]Shorten{},
	}
}

func (s *inmemService) PostShorten(ctx context.Context, p Shorten) (shorten ShortenResp, err error) {
	//s.mtx.Lock()
	//defer s.mtx.Unlock()
	//if _, ok := s.m[p.URL]; ok {
	//	return shorten, ErrAlreadyExists // POST = create, don't overwrite
	//}
	//s.m[p.URL] = p
	r2 := cache.CreateClient(1)
	defer r2.Close()

	ip := ctx.Value("ip").(string)
	if ip == "" {
		fmt.Println("could not get ip")
	}
	val, err := r2.Get(cache.Ctx, ip).Result()
	limit, _ := r2.TTL(cache.Ctx, ip).Result()

	if err == redis.Nil {
		_ = r2.Set(cache.Ctx, ip, os.Getenv("API_QUOTA"), 30*60*time.Second).Err()
	} else if err == nil {
		valInt, _ := strconv.Atoi(val)
		if valInt <= 0 {
			return shorten, errors.New("Rate limit exceeded, rate_limit_reset: " + fmt.Sprint(limit/time.Nanosecond/time.Minute))
		}
	}

	// check if the input is an actual URL

	if !govalidator.IsURL(p.URL) {
		return shorten, errors.New("Invalid URL")
	}

	// check for domain error
	if !helpers.RemoveDomainError(p.URL) {
		return shorten, errors.New("Can't do that: Domain error")
	}

	// enforce HTTPS, SSL
	p.URL = helpers.EnforceHTTP(p.URL)

	var id string
	if p.CustomShort != "" {
		id = p.CustomShort
	} else {
		id = helpers.Base62Encode(rand.Uint64())
	}

	r := cache.CreateClient(0)
	defer r.Close()

	val, _ = r.Get(cache.Ctx, id).Result()

	if val != "" {
		return shorten, errors.New("URL Custom short is already in use")
	}

	if p.Expiry == 0 {
		p.Expiry = 24
	}

	err = r.Set(cache.Ctx, id, p.URL, p.Expiry*3600*time.Second).Err()
	if err != nil {
		return shorten, errors.New("Unable to connect to server")
	}

	defaultAPIQuotaStr := os.Getenv("API_QUOTA")

	defaultApiQuota, _ := strconv.Atoi(defaultAPIQuotaStr)
	resp := ShortenResp{
		URL:             p.URL,
		CustomShort:     "",
		Expiry:          p.Expiry,
		XRateRemaining:  defaultApiQuota,
		XRateLimitReset: 30,
	}

	remainingQuota, err := r2.Decr(cache.Ctx, ip).Result()
	if err != nil {
		return shorten, errors.New("Unable to connect to server")
	}

	resp.XRateRemaining = int(remainingQuota)
	resp.XRateRemaining = int(limit / time.Nanosecond / time.Minute)

	resp.CustomShort = os.Getenv("DOMAIN") + "/" + id
	return resp, nil
}

func (s *inmemService) Resolve(ctx context.Context, url string) (string, error) {
	r := cache.CreateClient(0)
	defer r.Close()

	value, err := r.Get(cache.Ctx, url).Result()
	if err == redis.Nil {
		return "", errors.New("short-url not found in db")
	} else if err != nil {
		return "", errors.New("Internal error")
	}

	rInr := cache.CreateClient(1)
	defer rInr.Close()

	_ = rInr.Incr(cache.Ctx, "counter")

	return value, err
}
