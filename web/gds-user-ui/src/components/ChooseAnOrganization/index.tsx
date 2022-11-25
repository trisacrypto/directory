import {
  Grid,
  GridItem,
  HStack,
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalOverlay,
  Stack,
  StackDivider,
  Text
} from '@chakra-ui/react';
import { Account } from 'components/Account';
import AddNewVaspModal from 'components/AddNewVaspModal';
import InputFormControl from 'components/ui/InputFormControl';
import SelectFormControl from 'components/ui/SelectFormControl';
import { useOrganizationListQuery } from 'modules/dashboard/organization/useOrganizationListQuery';

const OPTIONS = [
  { label: 'Newest registrations', value: 'NEWEST_REGISTRATIONS' },
  { label: 'Most recently logged in', value: 'MOST_RECENTLY_LOGGED_IN' },
  { label: 'Alphabetical', value: 'ALPHABETICAL' }
];

export type ChooseAnAccountProps = {
  isOpen: boolean;
  onClose: () => void;
};

function ChooseAnOrganization({ isOpen, onClose }: ChooseAnAccountProps) {
  const { organizations } = useOrganizationListQuery();

  return (
    <>
      <Modal blockScrollOnMount={false} isOpen={isOpen} onClose={onClose} size="full">
        <ModalOverlay />
        <ModalContent pt="10vh">
          <ModalCloseButton />
          <ModalBody>
            <Stack spacing={3} maxW="700px" mx="auto">
              <div>
                <HStack width="100%" justifyContent="end">
                  <AddNewVaspModal />
                </HStack>
                <Text fontWeight={700}>Select an VASP from the Managed VASP List</Text>
                <Grid templateColumns="repeat(5, 1fr)" gap={4}>
                  <GridItem colSpan={3}>
                    <InputFormControl controlId="search" />
                  </GridItem>
                  <GridItem colSpan={2}>
                    <SelectFormControl controlId="filter" options={OPTIONS} />
                  </GridItem>
                </Grid>
              </div>
              <Stack divider={<StackDivider borderColor="#D9D9D9" />} p={2}>
                {organizations?.map((organization) => (
                  <Account
                    key={organization.id}
                    id={organization.id}
                    name={organization?.name}
                    domain={organization?.domain}
                  />
                ))}
              </Stack>
            </Stack>
          </ModalBody>
        </ModalContent>
      </Modal>
    </>
  );
}

export default ChooseAnOrganization;
