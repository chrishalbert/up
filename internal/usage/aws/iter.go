// Copyright 2023 Upbound Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aws

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"

	usagetime "github.com/upbound/up/internal/usage/time"
)

// ListObjectsV2InputIterator iterates through a *s3.ListObjectsV2Input for each
// window of time in a time range. Must be initialized with
// NewListObjectsV2InputIterator().
type ListObjectsV2InputIterator struct {
	Bucket  string
	Account string
	Iter    *usagetime.WindowIterator
}

// NewListObjectsV2InputIterator returns an initialized *ListObjectsV2InputIterator.
func NewListObjectsV2InputIterator(bucket string, account string, tr usagetime.Range, window time.Duration) (*ListObjectsV2InputIterator, error) {
	iter, err := usagetime.NewWindowIterator(tr, window)
	if err != nil {
		return nil, err
	}
	return &ListObjectsV2InputIterator{
		Bucket:  bucket,
		Account: account,
		Iter:    iter,
	}, nil
}

// More returns true if Next() has more to return.
func (i *ListObjectsV2InputIterator) More() bool {
	return i.Iter.More()
}

// Next returns a *s3.ListObjectsV2Input covering the next window of time, as
// well as a time range marking the window.
func (i *ListObjectsV2InputIterator) Next() (*s3.ListObjectsV2Input, usagetime.Range, error) {
	window, err := i.Iter.Next()
	if err != nil {
		return nil, usagetime.Range{}, err
	}
	// TODO(branden): Return []ListObjectsV2Input covering the complete window.
	return &s3.ListObjectsV2Input{
		Bucket: aws.String(i.Bucket),
		Prefix: aws.String(fmt.Sprintf(
			"account=%s/date=%s/hour=%02d/",
			i.Account,
			usagetime.FormatDateUTC(window.Start),
			window.Start.Hour(),
		)),
	}, window, nil
}
