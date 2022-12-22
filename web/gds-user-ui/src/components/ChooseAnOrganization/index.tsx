import { HStack, IconButton, Stack, StackDivider, Text, VStack } from '@chakra-ui/react';
import { Account } from 'components/Account';
import AddNewVaspModal from 'components/AddNewVaspModal/AddNewVaspModal';
import { Trans } from '@lingui/macro';
import { useOrganizationListQuery } from 'modules/dashboard/organization/useOrganizationListQuery';
import { userSelector } from 'modules/auth/login/user.slice';
import { useSelector } from 'react-redux';
import { GrClose } from 'react-icons/gr';
import { Link } from 'react-router-dom';

function ChooseAnOrganization() {
  const { organizations } = useOrganizationListQuery();
  const { user } = useSelector(userSelector);

  return (
    <VStack spacing={3} mx="auto" h="100vh" pt="10vh" position="relative">
      <IconButton
        as={Link}
        to="/"
        icon={<GrClose />}
        aria-label="Get back to dashboard"
        position="absolute"
        top={5}
        right={5}
        variant="ghost"
        title="Get back to dashboard"
      />
      <Stack maxW="700px" mx="auto">
        <div>
          <HStack width="100%" justifyContent="end">
            <AddNewVaspModal />
          </HStack>
          <Text fontWeight={700}>
            <Trans>Select a VASP from the Managed VASP List</Trans>
          </Text>
        </div>
        <Stack divider={<StackDivider borderColor="#D9D9D9" />} p={2}>
          {organizations && organizations?.length > 0 ? (
            organizations?.map((organization) => (
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
    </VStack>
  );
}

export default ChooseAnOrganization;
