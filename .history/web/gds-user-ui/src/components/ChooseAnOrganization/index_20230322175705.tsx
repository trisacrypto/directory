import { HStack, IconButton, Stack, StackDivider, Text, VStack, css } from '@chakra-ui/react';
import { Account } from 'components/Account';
import AddNewVaspModal from 'components/AddNewVaspModal/AddNewVaspModal';
import { Trans } from '@lingui/macro';
import { useOrganizationListQuery } from 'modules/dashboard/organization/useOrganizationListQuery';
import { userSelector } from 'modules/auth/login/user.slice';
import { useSelector } from 'react-redux';
import { GrClose } from 'react-icons/gr';
import { useNavigate } from 'react-router-dom';
import React from 'react';

// import { TransparentBackground } from 'components/TransparentBackground';
function ChooseAnOrganization() {
  const { organizations } = useOrganizationListQuery();
  console.log('[organizations list]', organizations);
  const { user } = useSelector(userSelector);
  const navigate = useNavigate();
  const handleBack = (e: React.MouseEvent<HTMLElement>) => {
    e.preventDefault();
    navigate(-1);
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
        width={'50%'}
        mx="auto"
        overflowY={'scroll'}
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
            {organizations && organizations?.organizations?.length > 0 ? (
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
    </VStack>
  );
}

export default ChooseAnOrganization;
