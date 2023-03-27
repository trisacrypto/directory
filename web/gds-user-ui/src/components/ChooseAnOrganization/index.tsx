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
import Loader from 'components/Loader';

// import { TransparentBackground } from 'components/TransparentBackground';
function ChooseAnOrganization() {
  const [currentPage, setCurrentPage] = useState<number>(1);
  // const [prevPage, setPrevPage] = useState(0);
  const [orgList, setOrgList] = useState<any>([]);
  const [wasLastList] = useState(false);
  const { organizations, getAllOrganizations, wasOrganizationFetched, isFetching } =
    useOrganizationListQuery(currentPage);

  const listInnerRef = useRef<any>();
  const { user } = useSelector(userSelector);
  const navigate = useNavigate();
  const handleBack = (e: React.MouseEvent<HTMLElement>) => {
    e.preventDefault();
    navigate(-1);
  };

  if (wasOrganizationFetched && orgList?.length === 0) {
    setOrgList(organizations?.organizations);
  }

  const NextPage = () => {
    setCurrentPage(currentPage + 1);
  };

  const PreviousPage = () => {
    setCurrentPage(currentPage - 1);
  };

  useEffect(() => {
    if (currentPage >= 1) {
      getAllOrganizations();
    }
  }, [currentPage, getAllOrganizations]);

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  // useEffect(() => {
  //   if (prevPage !== currentPage) {
  //     setPrevPage(currentPage + 1);
  //     if (organizations && organizations.organizations.length === 0) {
  //       setOrgList([...orgList, ...organizations.organizations]);
  //     }
  //   }
  // }, [currentPage, organizations, orgList, prevPage]);

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  // const onScroll = () => {
  //   if (listInnerRef.current) {
  //     const { scrollTop, scrollHeight, clientHeight } = listInnerRef.current;
  //     // console.log('[scrollTop]', scrollTop);
  //     // console.log('[scrollHeight]', scrollHeight);
  //     // console.log('[clientHeight]', clientHeight);
  //     // console.log('[scrollTop + clientHeight]', scrollTop + clientHeight);

  //     if (scrollTop + clientHeight === scrollHeight) {
  //       setCurrentPage(currentPage + 1);
  //     }
  //   }
  // };

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
      <div>
        <HStack width="100%" justify={'space-between'} spacing={20}>
          <Text fontWeight={700}>
            <Trans>Select a VASP from the Managed VASP List</Trans>
          </Text>

          <AddNewVaspModal />
        </HStack>
      </div>
      <Stack
        // onScroll={onScroll}
        ref={listInnerRef}
        width={'50%'}
        mx="auto"
        overflowY={'auto'}
        css={css({
          boxShadow: 'inset 0 -2px 0 rgba(0, 0, 0, 0.1)',
          border: '0 none'
        })}>
        <Stack>
          {isFetching && <Loader h="50vh" />}
          <Stack divider={<StackDivider borderColor="#D9D9D9" />} p={2}>
            {organizations?.organizations?.length > 0 ? (
              organizations?.organizations?.map((organization: any) => (
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
      <HStack>
        <Button onClick={PreviousPage} disabled={wasLastList}>
          Previous
        </Button>
        <Button onClick={NextPage} disabled={wasLastList}>
          Next
        </Button>
      </HStack>
    </VStack>
  );
}

export default ChooseAnOrganization;
