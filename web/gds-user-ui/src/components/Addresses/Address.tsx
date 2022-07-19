import { HStack, Box, Text } from '@chakra-ui/react';
import { Trans } from '@lingui/react';
import AddressForm from 'components/AddressForm';
import DeleteButton from 'components/ui/DeleteButton';

type AddressProps = {
  field: Record<'id', string>;
  index: number;
  onDelete: ((e: unknown) => void) | undefined;
};

function Address({ field, index, onDelete }: AddressProps) {
  return (
    <HStack key={field.id} width="100%" spacing={4} data-testid="address-row">
      <Box flex={1}>
        <Text>
          <Trans id="Address">Address</Trans> {index + 1}
        </Text>
        <AddressForm rowIndex={index} name={'entity.geographic_addresses'} />
      </Box>
      <Box alignSelf="flex-end" w={10} pb="25.1px">
        {index > 0 && (
          <DeleteButton
            onDelete={onDelete}
            tooltip={{
              label: <Trans id="Delete the address line">Delete the address line</Trans>
            }}
          />
        )}
      </Box>
    </HStack>
  );
}

export default Address;
