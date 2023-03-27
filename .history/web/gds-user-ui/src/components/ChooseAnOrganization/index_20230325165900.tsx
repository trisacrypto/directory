import {
  HStack,
  IconButton,
  Stack,
  StackDivider,
  Text,
  VStack,
  css,
  Button
} from '@chakra-ui/react';
import { Account } from 'components/Account';
import AddNewVaspModal from 'components/AddNewVaspModal/AddNewVaspModal';
import { Trans } from '@lingui/macro';
import { useOrganizationListQuery } from 'modules/dashboard/organization/useOrganizationListQuery';
import { userSelector } from 'modules/auth/login/user.slice';
import { useSelector } from 'react-redux';
import { GrClose } from 'react-icons/gr';
import { useNavigate } from 'react-router-dom';
import React, { useRef, useState, useEffect } from 'react';

// import { TransparentBackground } from 'components/TransparentBackground';
function ChooseAnOrganization() {
  const [currentPage, setCurrentPage] = useState(1);
  const [prevPage, setPrevPage] = useState(0);
  const [orgList, setOrgList] = useState<any>([]); // storing list
  const [wasLastList] = useState(false);
  const { organizations, getAllOrganizations } = useOrganizationListQuery(currentPage);

  const listInnerRef = useRef<any>();
  console.log('[organizations list]', organizations);
  const { user } = useSelector(userSelector);
  const navigate = useNavigate();
  const handleBack = (e: React.MouseEvent<HTMLElement>) => {
    e.preventDefault();
    navigate(-1);
  };

  const fetchMore = () => {
    setCurrentPage(currentPage + 1);
    getAllOrganizations();
  };

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  useEffect(() => {
    if (prevPage !== currentPage) {
      setPrevPage(currentPage + 1);
      if (organizations && organizations.organizations.length === 0) {
        setOrgList([...orgList, ...organizations.organizations]);
      }
    }
  }, [currentPage, organizations, orgList, prevPage]);

  useEffect(() => {
    if (organizations && organizations?.organizations.length === 0 && currentPage === 1) {
      setOrgList(organizations.organizations);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [organizations]);

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  const onScroll = () => {
    if (listInnerRef.current) {
      const { scrollTop, scrollHeight, clientHeight } = listInnerRef.current;
      if (scrollTop + clientHeight === scrollHeight) {
        setCurrentPage(currentPage + 1);
      }
    }
  };

  return (
    <VStack
      position="absolute"
      top="0"
      left="0"
      w="100%"
      h="100%"
      bg={'rgba(255,255,255,255)'}
      spacing={3}
      width={'full'}
      mx="auto"
      pt="10vh">
      <IconButton
        onClick={handleBack}
        icon={<GrClose />}
        aria-label="Get back to dashboard"
        position="absolute"
        top={5}
        right={5}
        variant="ghost"
        title="Get back to dashboard"
      />
      <Stack
        onScroll={onScroll}
        ref={listInnerRef}
        overflow="hidden"
        width={'50%'}
        mx="auto"
        // overflowY={'scroll'}
        css={css({
          boxShadow: 'inset 0 -2px 0 rgba(0, 0, 0, 0.1)',
          border: '0 none'
        })}>
        <div>
          <HStack width="100%" justifyContent="end">
            <AddNewVaspModal />
          </HStack>
          <Text fontWeight={700}>
            <Trans>Select a VASP from the Managed VASP List</Trans>
          </Text>
        </div>

        <Stack>
          <Stack divider={<StackDivider borderColor="#D9D9D9" />} p={2}>
            {orgList.length > 0 ? (
              orgList.map((organization: any) => (
                <Account
                  key={organization.id}
                  id={organization.id}
                  name={organization?.name}
                  domain={organization?.domain}
                  isCurrent={organization.id === user?.vasp?.id}
                />
              ))
            ) : (
              <Text>No VASPs found</Text>
            )}
          </Stack>
        </Stack>
      </Stack>
      <Button onClick={fetchMore} disabled={wasLastList}>
        Load more
      </Button>
    </VStack>
  );
}

export default ChooseAnOrganization;
