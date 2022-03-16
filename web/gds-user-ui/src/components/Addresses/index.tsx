import { DeleteIcon } from '@chakra-ui/icons';
import { Box, Button, Heading, HStack, Icon, Text, Tooltip, VStack } from '@chakra-ui/react';
import DeleteButton from 'components/ui/DeleteButton';
import FormButton from 'components/ui/FormButton';
import FormLayout from 'layouts/FormLayout';
import AddressForm from '../AddressForm';

type AddressesPropsProps = {};
const Addresses: React.FC<AddressesPropsProps> = () => {
  return (
    <FormLayout>
      <Heading size="md">Addresses</Heading>
      <Text size="sm">Enter at least one geographic address</Text>
      <VStack width="100%" align="start" spacing={4}>
        <HStack width="100%" spacing={4}>
          <Box flex={1}>
            <AddressForm />
          </Box>
          <Box alignSelf="flex-end" w={10} pb="25.1px">
            <DeleteButton tooltip={{ label: 'Delete the address line' }} />
          </Box>
        </HStack>
        <Box>
          <FormButton borderRadius="5px">Add Address</FormButton>
        </Box>
      </VStack>
    </FormLayout>
  );
};

export default Addresses;
