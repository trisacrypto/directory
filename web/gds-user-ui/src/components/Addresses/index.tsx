import { Box, Button, Heading, HStack, Text, Tooltip, VStack } from '@chakra-ui/react';
import DeleteButton from 'components/ui/DeleteButton';
import FormButton from 'components/ui/FormButton';
import FormLayout from 'layouts/FormLayout';
import React from 'react';
import { useFieldArray, useFormContext } from 'react-hook-form';
import AddressForm from '../AddressForm';

type AddressesPropsProps = {};
const Addresses: React.FC<AddressesPropsProps> = () => {
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
      <Heading size="md">Addresses</Heading>
      <Text size="sm">Enter at least one geographic address</Text>
      <VStack width="100%" align="start" spacing={10}>
        {fields.map((field, index) => {
          return (
            <HStack key={field.id} width="100%" spacing={4}>
              <Box flex={1}>
                <Text>Address {index + 1}</Text>
                <AddressForm
                  rowIndex={index}
                  name={'entity.geographic_addresses'}
                  register={register}
                  control={control}
                />
              </Box>
              <Box alignSelf="flex-end" w={10} pb="25.1px">
                {index ? (
                  <DeleteButton
                    onDelete={() => remove(index)}
                    tooltip={{ label: 'Delete the address line' }}
                  />
                ) : null}
              </Box>
            </HStack>
          );
        })}
        <Box>
          <FormButton onClick={handleAddressClick} borderRadius="5px">
            Add Address
          </FormButton>
        </Box>
      </VStack>
    </FormLayout>
  );
};

export default Addresses;
