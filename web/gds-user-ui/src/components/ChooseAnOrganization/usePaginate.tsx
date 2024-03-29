import { useState, useEffect } from 'react';
import { useOrganizationListQuery } from 'modules/dashboard/organization/useOrganizationListQuery';

export const useOrganizationPagination = (searchQuery?: string) => {
  const [currentPage, setCurrentPage] = useState<number>(1);
  const [wasLastPage, setWasLastPage] = useState<boolean>(false);
  const { organizations, getAllOrganizations, isFetching } = useOrganizationListQuery({
    page: currentPage,
    name: searchQuery
  });
  const { count, page_size, page } = organizations || {};

  const NextPage = () => {
    setCurrentPage(currentPage + 1);
  };

  const PreviousPage = () => {
    setCurrentPage(currentPage - 1);
  };

  useEffect(() => {
    if (searchQuery || searchQuery === '') {
      getAllOrganizations();
    }
  }, [searchQuery, getAllOrganizations]);

  useEffect(() => {
    if (currentPage !== 1) {
      getAllOrganizations();
    }
  }, [currentPage, getAllOrganizations, searchQuery]);

  useEffect(() => {
    if (page && page_size && count && page * page_size >= count) {
      setWasLastPage(true);
    }
    return () => {
      setWasLastPage(false);
    };
  }, [page, page_size, count]);

  return {
    NextPage,
    PreviousPage,
    currentPage,
    wasLastPage,
    isFetching,
    organizations: organizations?.organizations || []
  };
};
