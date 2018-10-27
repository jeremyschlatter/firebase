package firebase

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/option"

	firebaseDB "github.com/jeremyschlatter/firebase/db"
)

var (
	projectID   string
	credentials option.ClientOption
)

func init() {
	// Parse admin SDK credentials.

	adminSDK := "firebase-admin-sdk-key.json"
	keyBytes, err := ioutil.ReadFile(adminSDK)

	if os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, `
These tests require a live, empty Firebase project.

Before running these tests, you should:

    1. Create a new, empty Firebase project.

       This test will WIPE THE DATABASE of your project, so don't reuse an existing project unless you're ok with losing all of its data.

    2. Download service account credentials for your project.

       Use Google to find out how to do this, or try finding it in the Firebase console:

	       Project settings > Service accounts > Firebase Admin SDK

    3. Put the credentials in the same directory as this test, with the filename %q.

`, adminSDK)
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	var data struct {
		ProjectID string `json:"project_id"`
	}
	err = json.Unmarshal(keyBytes, &data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Parsing %v: %v\n", adminSDK, err)
		os.Exit(1)
	}

	projectID = data.ProjectID
	credentials = option.WithCredentialsJSON(keyBytes)
}

func fakeDB(t *testing.T) firebaseDB.Client {
	db, _ := NewFake().Database(context.Background())
	return db
}

func realDB(t *testing.T) firebaseDB.Client {
	ctx := context.Background()
	app, err := NewApp(
		ctx,
		&Config{
			DatabaseURL: "https://" + projectID + ".firebaseio.com",
			ProjectID:   projectID,
		},
		credentials,
	)
	require.Nil(t, err)

	db, err := app.Database(ctx)
	require.Nil(t, err)

	require.Nil(t, db.NewRef("/").Delete(ctx))

	return db
}

func TestFake(t *testing.T) {
	ctx := context.Background()
	db := fakeDB(t)

	original := "hello"
	assert.Nil(t, db.NewRef("/test").Set(ctx, original))

	var got interface{}
	assert.Nil(t, db.NewRef("/test").Get(ctx, &got))
	assert.Equal(t, original, got)
}

func TestReal(t *testing.T) {
	ctx := context.Background()

	db := realDB(t)

	original := "hello"
	assert.Nil(t, db.NewRef("/test").Set(ctx, original))

	var got interface{}
	assert.Nil(t, db.NewRef("/test").Get(ctx, &got))
	assert.Equal(t, original, got)
}

func TestEquivalence(t *testing.T) {
	ctx := context.Background()

	interact := func(db firebaseDB.Client) (results []string) {
		check := func(ref firebaseDB.Ref) {
			var i interface{}
			assert.Nil(t, ref.Get(ctx, &i))
			b, err := json.Marshal(i)
			assert.Nil(t, err)
			results = append(results, string(b))
		}

		root := db.NewRef("/")
		check(root)

		db.NewRef("/").Set(ctx, "hello")
		check(root)

		root.Set(ctx, 1)
		check(root)

		var data struct {
			Foo string
			Bar int
			Sub struct {
				Quux int
			}
		}
		db.NewRef("/").Set(ctx, data)
		check(root)

		db.NewRef("/a/b/c").Set(ctx, 12)
		check(db.NewRef("/a").Child("b").Child("c"))

		return results
	}

	assert.Equal(t, interact(realDB(t)), interact(fakeDB(t)))
}
