/*
 * Copyright (c) 2023 Zander Schwid & Co. LLC.
 * SPDX-License-Identifier: BUSL-1.1
 */

package sprint

import (
	"context"
	"github.com/codeallergy/glue"
	"github.com/sprintframework/sprintpb"
	"github.com/keyvalstore/store"
	"github.com/codeallergy/uuid"
	htmlTemplate "html/template"
	"io"
	"reflect"
	textTemplate "text/template"
	"time"
)


var ResourceServiceClass = reflect.TypeOf((*ResourceService)(nil)).Elem()

type ResourceService interface {

	/*
	Gets resource by name
	 */
	GetResource(name string) ([]byte, error)

	/*
	Gets text template resource by name
	 */
	TextTemplate(name string) (*textTemplate.Template, error)

	/*
	Gets html template resource by name
	 */
	HtmlTemplate(name string) (*htmlTemplate.Template, error)

	/*
	Gets using licences of imported modules
	 */
	GetLicenses(name string) (string, error)

	/*
	Gets open api swagger JSON files for resource source
	 */
	GetOpenAPI(source string) string
}


var ConfigRepositoryClass = reflect.TypeOf((*ConfigRepository)(nil)).Elem()

type ConfigRepository interface {
	glue.DisposableBean
	glue.PropertyResolver

	/**
	Gets property value the property name (key) if found or default value.

	If property not found in storage then function will return empty string with no error.

	In case of issue function will return error.
	*/

	Get(key string) (string, error)

	/**
	Enumerates all properties that start with provided prefix.
	Prefix could be an empty string, in this case all properties will be enumerated.

	On each call callback function should return true to continue enumeration.

	In case of issue function will return error.
	*/

	EnumerateAll(prefix string, cb func(key, value string) bool) error

	/**
	Sets specific string property with key.
	If value is empty string, then the property would be removed from config storage.
	All properties are stored in string values on backend.

	In case of issue function will return error.
	*/

	Set(key, value string) error

	/**
	Watch updates with prefix on backend system during specific active context.

	On each call callback function should return true to continue watching on changes.

	If property deleted, then watching value will be empty for the specific key.

	In case of issue function will return error.

	Use Application as context.
	*/

	Watch(context context.Context, prefix string, cb func(key, value string) bool) (context.CancelFunc, error)

	/**
	Gets backend using for storing properties
	*/

	Backend() store.DataStore

	/**
	Sets backend using for storing properties
	*/
	SetBackend(storage store.DataStore)

}

var AutoupdateServiceClass = reflect.TypeOf((*AutoupdateService)(nil)).Elem()

type AutoupdateService interface {
	glue.InitializingBean
	glue.DisposableBean

	/**
	Suspend the job
	 */

	Freeze(jobName string) int64

	/**
	Revoke the job
	 */

	Unfreeze(handle int64)

	/**
	Get list of jobs
	 */

	FreezeJobs() map[int64]string

}

var NodeServiceClass = reflect.TypeOf((*NodeService)(nil)).Elem()

type NodeService interface {
	glue.InitializingBean
	Component

	/**
	Returns the node id unique number. Usually random number defined on the first startup.
	 */

	NodeId() uint64

	/**
	Returns node is in hex format as a string. Could be used in Raft and Gossip protocols as unique node number.
	 */

	NodeIdHex() string

	/**
	Issue random time based UUID number having node id, current timestamp and random number.
	Perfect candidate to identity ordered events in distributed systems.
	 */

	Issue() uuid.UUID

	/**
	Parses UUID issues by this service. Clock usually the random number during generation.
	 */

	Parse(uuid.UUID) (timestampMillis int64, nodeId int64, clock int)
}

type StorageConsoleStream interface {

	/**
	Grpc streaming interface to send multiple responses in bi-directional steams.
	 */

	Send(*sprintpb.StorageConsoleResponse) error

	/**
	Grpc streaming interface to receive multiple requests in bi-directional steams.
	 */

	Recv() (*sprintpb.StorageConsoleRequest, error)
}

var StorageServiceClass = reflect.TypeOf((*StorageService)(nil)).Elem()

type StorageService interface {
	glue.InitializingBean

	/**
	Executes query on storage.
	 */

	ExecuteQuery(name, query string, cb func(string) bool) error

	/**
	Executes command on storage.
	 */

	ExecuteCommand(cmd string, args []string) (string, error)

	/**
	Starts bi-directional console with embedded storage through gRPC streams.
	 */

	Console(stream StorageConsoleStream) error

	/**
	Starts bi-directional console with embedded storage.
	 */

	LocalConsole(writer io.StringWriter, errWriter io.StringWriter) error

}

var JobServiceClass = reflect.TypeOf((*JobService)(nil)).Elem()

type JobInfo struct {

	/**
	Unique job name
	 */

	Name         string

	/**
	Schedule format, when and how often.
	 */

	Schedule     string

	/**
	Client function
	 */

	ExecutionFn  func(context.Context) error
}

type JobService interface {

	/**
	List all scheduled and running jobs
	 */

	ListJobs() ([]string, error)

	/**
	Add background job
	 */

	AddJob(*JobInfo) error

	/**
	Cancels job by name
	 */

	CancelJob(name string) error

	/**
	Runs the job by name
	 */

	RunJob(ctx context.Context, name string) error

	/**
	Executes command on job service
	 */

	ExecuteCommand(cmd string, args []string) (string, error)

}

var AuthenticationServiceClass = reflect.TypeOf((*AuthorizationMiddleware)(nil)).Elem()

type AuthorizedUser struct {

	/**
	User id or unique username of the authorized user
	 */

	Username   string

	/**
	Roles assigned to the user
	 */

	Roles      map[string]bool

	/**
	Additional permissions and ACL list for the user
	 */

	Context    map[string]string

	/**
	Expiration of the token in milliseconds
	 */

	ExpiresAt  int64

	/**
	JWT auth token
	 */
	Token      string
}

type AuthorizationMiddleware interface {
	glue.InitializingBean

	/**
	Authenticates user by using metadata from gRPC context.
	Places the AuthorizedUser object in to context, returns new combined context.
	Runs from middleware automatically by gRPC Server on each request.
	 */

	Authenticate(ctx context.Context) (context.Context, error)

	/**
	Gets AuthorizedUser from the Context. In case of missing of the object calls Authenticate.
	 */

	GetUser(ctx context.Context) (*AuthorizedUser, bool)

	/**
	Checks if context has user role.
	 */

	HasUserRole(ctx context.Context, role string) bool

	/**
	Gets additional ACL list from the context.
	 */

	UserContext(ctx context.Context, name string) (string, bool)

	/**
	Generates JWT token for AuthorizedUser. Token field would be reset.
	 */

	GenerateToken(user *AuthorizedUser) (string, error)

	/**
	Parses token and returns AuthorizedUser object.
	 */

	ParseToken(token string) (*AuthorizedUser, error)

	/**
	Adds token to invalidate list. The next login would not be possible with it.
	 */

	InvalidateToken(token string)

}

var MailServiceClass = reflect.TypeOf((*MailService)(nil)).Elem()

type Mail struct {

	/**
	Sender email address with optional name.
	 */

	Sender  string

	/**
	Non-empty list of recipient addressed with optional names.
	 */

	Recipients []string

	/**
	Subject of the message.
	 */

	Subject string

	/**
	Text body of the message.
	 */

	TextTemplate string

	/**
	Optional HTML body of the message.
	 */

	HtmlTemplate string   // optional

	/**
	Body of the message.
	 */

	Data interface{}

	/**
	Optional attachments to the message.
	 */
	Attachments []string  // optional
}

type MailService interface {
	glue.NamedBean

	/**
	Sends the message.
	 */

	SendMail(mail *Mail, timeout time.Duration, async bool) error

}

