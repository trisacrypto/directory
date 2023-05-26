import { VStack, Box, Button } from '@chakra-ui/react';
import { Trans } from '@lingui/react';
import { useFormContext, useFieldArray } from 'react-hook-form';
import Address from './Address';

function AddressList() {
  const { control } = useFormContext();
  const { fields, remove, append } = useFieldArray({
    control,
    name: 'entity.geographic_addresses'
  });

  const handleAddressClick = () => {
    append({
      address_type: '',
      address_line: ['', '', ''],
      country: ''
    });
  };

  return (
    <VStack width="100%" align="start" spacing={10} data-testid="legal-adress">
      {fields.map((field, index) => {
        return (
          <Address key={field.id} field={field} index={index} onDelete={() => remove(index)} />
        );
      })}
      <Box>
        <Button onClick={handleAddressClick}>
          <Trans id="Add Address">Add Address</Trans>
        </Button>
      </Box>
    </VStack>
  );
}

export default AddressList;
