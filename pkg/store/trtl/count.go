package trtl

import (
	"context"

	"github.com/trisacrypto/directory/pkg/trtl/pb/v1"
	"github.com/trisacrypto/directory/pkg/utils/wire"
)

func (s *Store) CountVASPs(ctx context.Context) (_ uint64, err error) {
	var reply *pb.CountReply
	if reply, err = s.client.Count(ctx, &pb.CountRequest{Namespace: wire.NamespaceVASPs}); err != nil {
		return 0, err
	}
	return reply.Objects, nil
}

func (s *Store) CountCertReqs(ctx context.Context) (_ uint64, err error) {
	var reply *pb.CountReply
	if reply, err = s.client.Count(ctx, &pb.CountRequest{Namespace: wire.NamespaceCertReqs}); err != nil {
		return 0, err
	}
	return reply.Objects, nil
}

func (s *Store) CountCerts(ctx context.Context) (_ uint64, err error) {
	var reply *pb.CountReply
	if reply, err = s.client.Count(ctx, &pb.CountRequest{Namespace: wire.NamespaceCerts}); err != nil {
		return 0, err
	}
	return reply.Objects, nil
}

func (s *Store) CountAnnouncementMonths(ctx context.Context) (_ uint64, err error) {
	var reply *pb.CountReply
	if reply, err = s.client.Count(ctx, &pb.CountRequest{Namespace: wire.NamespaceAnnouncements}); err != nil {
		return 0, err
	}
	return reply.Objects, nil
}

func (s *Store) CountActivityMonth(ctx context.Context) (_ uint64, err error) {
	var reply *pb.CountReply
	if reply, err = s.client.Count(ctx, &pb.CountRequest{Namespace: wire.NamespaceActivities}); err != nil {
		return 0, err
	}
	return reply.Objects, nil
}

func (s *Store) CountOrganizations(ctx context.Context) (_ uint64, err error) {
	var reply *pb.CountReply
	if reply, err = s.client.Count(ctx, &pb.CountRequest{Namespace: wire.NamespaceOrganizations}); err != nil {
		return 0, err
	}
	return reply.Objects, nil
}

func (s *Store) CountContacts(ctx context.Context) (_ uint64, err error) {
	var reply *pb.CountReply
	if reply, err = s.client.Count(ctx, &pb.CountRequest{Namespace: wire.NamespaceContacts}); err != nil {
		return 0, err
	}
	return reply.Objects, nil
}

func (s *Store) CountEmails(ctx context.Context) (_ uint64, err error) {
	var reply *pb.CountReply
	if reply, err = s.client.Count(ctx, &pb.CountRequest{Namespace: wire.NamespaceEmails}); err != nil {
		return 0, err
	}
	return reply.Objects, nil
}
