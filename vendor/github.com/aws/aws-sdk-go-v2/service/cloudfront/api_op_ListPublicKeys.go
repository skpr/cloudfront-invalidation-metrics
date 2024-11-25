// Code generated by smithy-go-codegen DO NOT EDIT.

package cloudfront

import (
	"context"
	"fmt"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// List all public keys that have been added to CloudFront for this account.
func (c *Client) ListPublicKeys(ctx context.Context, params *ListPublicKeysInput, optFns ...func(*Options)) (*ListPublicKeysOutput, error) {
	if params == nil {
		params = &ListPublicKeysInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "ListPublicKeys", params, optFns, c.addOperationListPublicKeysMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*ListPublicKeysOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type ListPublicKeysInput struct {

	// Use this when paginating results to indicate where to begin in your list of
	// public keys. The results include public keys in the list that occur after the
	// marker. To get the next page of results, set the Marker to the value of the
	// NextMarker from the current page's response (which is also the ID of the last
	// public key on that page).
	Marker *string

	// The maximum number of public keys you want in the response body.
	MaxItems *int32

	noSmithyDocumentSerde
}

type ListPublicKeysOutput struct {

	// Returns a list of all public keys that have been added to CloudFront for this
	// account.
	PublicKeyList *types.PublicKeyList

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationListPublicKeysMiddlewares(stack *middleware.Stack, options Options) (err error) {
	if err := stack.Serialize.Add(&setOperationInputMiddleware{}, middleware.After); err != nil {
		return err
	}
	err = stack.Serialize.Add(&awsRestxml_serializeOpListPublicKeys{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsRestxml_deserializeOpListPublicKeys{}, middleware.After)
	if err != nil {
		return err
	}
	if err := addProtocolFinalizerMiddlewares(stack, options, "ListPublicKeys"); err != nil {
		return fmt.Errorf("add protocol finalizers: %v", err)
	}

	if err = addlegacyEndpointContextSetter(stack, options); err != nil {
		return err
	}
	if err = addSetLoggerMiddleware(stack, options); err != nil {
		return err
	}
	if err = addClientRequestID(stack); err != nil {
		return err
	}
	if err = addComputeContentLength(stack); err != nil {
		return err
	}
	if err = addResolveEndpointMiddleware(stack, options); err != nil {
		return err
	}
	if err = addComputePayloadSHA256(stack); err != nil {
		return err
	}
	if err = addRetry(stack, options); err != nil {
		return err
	}
	if err = addRawResponseToMetadata(stack); err != nil {
		return err
	}
	if err = addRecordResponseTiming(stack); err != nil {
		return err
	}
	if err = addSpanRetryLoop(stack, options); err != nil {
		return err
	}
	if err = addClientUserAgent(stack, options); err != nil {
		return err
	}
	if err = smithyhttp.AddErrorCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = addSetLegacyContextSigningOptionsMiddleware(stack); err != nil {
		return err
	}
	if err = addTimeOffsetBuild(stack, c); err != nil {
		return err
	}
	if err = addUserAgentRetryMode(stack, options); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opListPublicKeys(options.Region), middleware.Before); err != nil {
		return err
	}
	if err = addRecursionDetection(stack); err != nil {
		return err
	}
	if err = addRequestIDRetrieverMiddleware(stack); err != nil {
		return err
	}
	if err = addResponseErrorMiddleware(stack); err != nil {
		return err
	}
	if err = addRequestResponseLogging(stack, options); err != nil {
		return err
	}
	if err = addDisableHTTPSMiddleware(stack, options); err != nil {
		return err
	}
	if err = addSpanInitializeStart(stack); err != nil {
		return err
	}
	if err = addSpanInitializeEnd(stack); err != nil {
		return err
	}
	if err = addSpanBuildRequestStart(stack); err != nil {
		return err
	}
	if err = addSpanBuildRequestEnd(stack); err != nil {
		return err
	}
	return nil
}

// ListPublicKeysPaginatorOptions is the paginator options for ListPublicKeys
type ListPublicKeysPaginatorOptions struct {
	// The maximum number of public keys you want in the response body.
	Limit int32

	// Set to true if pagination should stop if the service returns a pagination token
	// that matches the most recent token provided to the service.
	StopOnDuplicateToken bool
}

// ListPublicKeysPaginator is a paginator for ListPublicKeys
type ListPublicKeysPaginator struct {
	options   ListPublicKeysPaginatorOptions
	client    ListPublicKeysAPIClient
	params    *ListPublicKeysInput
	nextToken *string
	firstPage bool
}

// NewListPublicKeysPaginator returns a new ListPublicKeysPaginator
func NewListPublicKeysPaginator(client ListPublicKeysAPIClient, params *ListPublicKeysInput, optFns ...func(*ListPublicKeysPaginatorOptions)) *ListPublicKeysPaginator {
	if params == nil {
		params = &ListPublicKeysInput{}
	}

	options := ListPublicKeysPaginatorOptions{}
	if params.MaxItems != nil {
		options.Limit = *params.MaxItems
	}

	for _, fn := range optFns {
		fn(&options)
	}

	return &ListPublicKeysPaginator{
		options:   options,
		client:    client,
		params:    params,
		firstPage: true,
		nextToken: params.Marker,
	}
}

// HasMorePages returns a boolean indicating whether more pages are available
func (p *ListPublicKeysPaginator) HasMorePages() bool {
	return p.firstPage || (p.nextToken != nil && len(*p.nextToken) != 0)
}

// NextPage retrieves the next ListPublicKeys page.
func (p *ListPublicKeysPaginator) NextPage(ctx context.Context, optFns ...func(*Options)) (*ListPublicKeysOutput, error) {
	if !p.HasMorePages() {
		return nil, fmt.Errorf("no more pages available")
	}

	params := *p.params
	params.Marker = p.nextToken

	var limit *int32
	if p.options.Limit > 0 {
		limit = &p.options.Limit
	}
	params.MaxItems = limit

	optFns = append([]func(*Options){
		addIsPaginatorUserAgent,
	}, optFns...)
	result, err := p.client.ListPublicKeys(ctx, &params, optFns...)
	if err != nil {
		return nil, err
	}
	p.firstPage = false

	prevToken := p.nextToken
	p.nextToken = nil
	if result.PublicKeyList != nil {
		p.nextToken = result.PublicKeyList.NextMarker
	}

	if p.options.StopOnDuplicateToken &&
		prevToken != nil &&
		p.nextToken != nil &&
		*prevToken == *p.nextToken {
		p.nextToken = nil
	}

	return result, nil
}

// ListPublicKeysAPIClient is a client that implements the ListPublicKeys
// operation.
type ListPublicKeysAPIClient interface {
	ListPublicKeys(context.Context, *ListPublicKeysInput, ...func(*Options)) (*ListPublicKeysOutput, error)
}

var _ ListPublicKeysAPIClient = (*Client)(nil)

func newServiceMetadataMiddleware_opListPublicKeys(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		OperationName: "ListPublicKeys",
	}
}
