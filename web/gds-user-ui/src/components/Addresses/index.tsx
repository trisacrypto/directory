import { Box, Button, Heading, HStack, Text, Tooltip, VStack } from '@chakra-ui/react';
import { Trans } from '@lingui/react';
import DeleteButton from 'components/ui/DeleteButton';
import FormLayout from 'layouts/FormLayout';
import React from 'react';
import { useFieldArray, useFormContext } from 'react-hook-form';
import AddressForm from '../AddressForm';

const Addresses: React.FC = () => {
  const { control, register } = useFormContext();
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
    <FormLayout>
      <Heading size="md">
        <Trans id="Addresses">Addresses</Trans>
      </Heading>
      <Text size="sm">
        <Trans id="At least one geographic address is required. Enter the primary geographic address of the organization. Organizations may enter additional addresses if operating in multiple jurisdictions.">
          At least one geographic address is required. Enter the primary geographic address of the
          organization. Organizations may enter additional addresses if operating in multiple
          jurisdictions.
        </Trans>
      </Text>
      <VStack width="100%" align="start" spacing={10}>
        {fields.map((field, index) => {
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
                    onDelete={() => remove(index)}
                    tooltip={{
                      label: <Trans id="Delete the address line">Delete the address line</Trans>
                    }}
                  />
                )}
              </Box>
            </HStack>
          );
        })}
        <Box>
          <Button onClick={handleAddressClick}>
            <Trans id="Add Address">Add Address</Trans>
          </Button>
        </Box>
      </VStack>
    </FormLayout>
  );
};

export default Addresses;
