import { useEffect } from 'react';
import { Grid, GridItem, VStack, Text } from '@chakra-ui/react';
import InputFormControl from 'components/ui/InputFormControl';
import SelectFormControl from 'components/ui/SelectFormControl';
import { addressTypeOptions } from 'constants/address';
import { getCountriesOptions } from 'constants/countries';
import { Control, Controller, useFormContext, UseFormRegister } from 'react-hook-form';
import _ from 'lodash';
import { getValueByPathname } from 'utils/utils';

type AddressFormProps = {
  control: Control;
  register: UseFormRegister<any>;
  name: string;
  rowIndex: number;
};

const AddressForm: React.FC<AddressFormProps> = ({ register, control, name, rowIndex }) => {
  const countries = getCountriesOptions();
  const addressTypes = addressTypeOptions();
  const {
    watch,
    formState: { errors },
    setValue
  } = useFormContext();

  const getFirstAddressType = watch('entity.geographic_addresses[0].address_type');

  useEffect(() => {
    if (!getFirstAddressType) {
      setValue(`entity.geographic_addresses[0].address_type`, 'ADDRESS_TYPE_CODE_BIZZ');
    }
  }, [getFirstAddressType]);

  return (
    <>
      <VStack spacing={3.5} align="start">
        <InputFormControl
          formHelperText="Address line 1 e.g. building name/number, street name"
          controlId={`${name}[${rowIndex}].address_line[0]`}
          isInvalid={!!getValueByPathname(errors, `${name}[${rowIndex}].address_line[0]`)}
          {...register(`${name}[${rowIndex}].address_line[0]`)}
        />

        <InputFormControl
          formHelperText="Address line 2 e.g. apartment or suite number"
          controlId="address_2"
          isInvalid={!!getValueByPathname(errors, `${name}[${rowIndex}].address_line[1]`)}
          {...register(`${name}[${rowIndex}].address_line[1]`)}
        />

        <InputFormControl
          formHelperText="Address line 3 e.g. city, province, postal code"
          controlId="address_3"
          isInvalid={!!getValueByPathname(errors, `${name}[${rowIndex}].address_line[2]`)}
          {...register(`${name}[${rowIndex}].address_line[2]`)}
        />

        <Grid templateColumns={{ base: '1fr', md: 'repeat(2, 1fr)' }} gap={6} width="100%">
          <GridItem>
            <Controller
              control={control}
              name={`${name}[${rowIndex}].country`}
              render={({ field }) => (
                <SelectFormControl
                  name={field.name}
                  ref={field.ref}
                  options={countries}
                  isInvalid={!!getValueByPathname(errors, `${name}[${rowIndex}].country`)}
                  value={countries.find((option) => option.value === field.value)}
                  onChange={(newValue: any) => field.onChange(newValue.value)}
                  formHelperText="Country"
                  controlId="country"
                />
              )}
            />
          </GridItem>
          <GridItem>
            <Controller
              control={control}
              name={`${name}[${rowIndex}].address_type`}
              render={({ field }) => (
                <SelectFormControl
                  name={field.name}
                  ref={field.ref}
                  isInvalid={!!getValueByPathname(errors, `${name}[${rowIndex}].address_type`)}
                  value={addressTypes.find((option) => option.value === field.value)}
                  onChange={(newValue: any) => field.onChange(newValue.value)}
                  options={addressTypes}
                  formHelperText="Address Type"
                  controlId="address_type"
                />
              )}
            />
          </GridItem>
        </Grid>
      </VStack>
    </>
  );
};

export default AddressForm;
