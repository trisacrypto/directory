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

import { userSelector } from 'modules/auth/login/user.slice';
import { useSelector } from 'react-redux';
import { GrClose } from 'react-icons/gr';
import { useNavigate } from 'react-router-dom';
import Loader from 'components/Loader';
import { useOrganizationPagination } from './usePaginate';
import { useState } from 'react';
import SearchVasp from './SearchVasp';
function ChooseAnOrganization() {
  const [searchQuery, setSearchQuery] = useState<string>('');
  const { NextPage, PreviousPage, currentPage, wasLastPage, isFetching, organizations } =
    useOrganizationPagination(searchQuery);

  const { user } = useSelector(userSelector);
  const navigate = useNavigate();

  const handleBack = (e: React.MouseEvent<HTMLElement>) => {
    e.preventDefault();
    navigate(-1);
  };

  const isFirstPage = currentPage === 1;

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
      <Text fontWeight={700}>
        <Trans>Select a VASP from the Managed VASP List</Trans>
      </Text>

      <Stack
        mx="auto"
        overflowY={'auto'}
        css={css({
          boxShadow: 'inset 0 -2px 0 rgba(0, 0, 0, 0.1)',
          border: '0 none'
        })}>
        <HStack justify={'space-between'}>
          <SearchVasp setSearchOrganization={setSearchQuery} />
          <AddNewVaspModal />
        </HStack>
        <Stack>
          {isFetching && <Loader h="50vh" />}
          <Stack divider={<StackDivider borderColor="#D9D9D9" />} p={2}>
            {organizations?.length > 0 ? (
              organizations?.map((organization: any) => (
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
        <Button onClick={PreviousPage} disabled={isFirstPage}>
          Previous
        </Button>
        <Button onClick={NextPage} disabled={!!wasLastPage}>
          Next
        </Button>
      </HStack>
    </VStack>
  );
}

export default ChooseAnOrganization;
