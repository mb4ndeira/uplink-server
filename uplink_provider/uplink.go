package uplink_provider

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"storj.io/uplink"
	"storj.io/uplink/edge"
)

type Server struct {
	UplinkServiceServer
}

func generateURL(ctx context.Context, access *uplink.Access, bucket string, key string, baseURL string) error {
	publicAccess := true

	var edgeConfig edge.Config

	// edgeConfig.AuthServiceAddress =
	// edgeConfig.CertificatePEM =

	credentials, err := edgeConfig.RegisterAccess(ctx, access, &edge.RegisterAccessOptions{Public: publicAccess})
	if err != nil {
		return err
	}

	url, err := edge.JoinShareURL(baseURL, credentials.AccessKeyID, bucket, key, nil)
	if err != nil {
		return err
	}

	fmt.Print(url)

	return nil
}

func (s *Server) Upload(stream UplinkService_UploadServer) error {
	ctx := context.Background()

	args, err := stream.Recv()
	if err != nil && err != io.EOF {
		return err
	}

	access, err := uplink.ParseAccess(args.AccessGrant)
	if err != nil {
		return fmt.Errorf("could not parse access grant: %w", err)
	}

	project, err := uplink.OpenProject(ctx, access)
	if err != nil {
		return fmt.Errorf("could not open project: %w", err)
	}
	defer project.Close()

	_, err = project.EnsureBucket(ctx, args.BucketName)
	if err != nil {
		return fmt.Errorf("could not ensure bucket: %w", err)
	}

	upload, err := project.UploadObject(ctx, args.BucketName, args.ObjectKey, nil)
	if err != nil {
		return fmt.Errorf("could not initiate upload: %w", err)
	}

	_, err = io.Copy(upload, bytes.NewReader(args.Data))
	if err != nil {
		_ = upload.Abort()
		return fmt.Errorf("could not upload data: %w", err)
	}

	err = upload.Commit()
	if err != nil {
		return fmt.Errorf("could not commit uploaded object: %w", err)
	}

	err = generateURL(ctx, access, args.BucketName, args.ObjectKey, "base-url")
	if err != nil {
		return fmt.Errorf("could not generate public URL: %w", err)
	}

	stream.SendAndClose(&UploadReturn{Response: &UploadResponse{Url: "fake-url"}})

	return nil
}
