import React, { useEffect } from 'react';
import { Grid, GridItem, VStack } from '@chakra-ui/react';
import InputFormControl from 'components/ui/InputFormControl';
import SelectFormControl from 'components/ui/SelectFormControl';
import { addressTypeEnum, addressTypeOptions } from 'constants/address';
import { getCountriesOptions } from 'constants/countries';
import { Controller, useFormContext } from 'react-hook-form';
import { getValueByPathname } from 'utils/utils';
import { t } from '@lingui/macro';

type AddressFormProps = {
  name: string;
  rowIndex: number;
};

const AddressForm: React.FC<AddressFormProps> = ({ name, rowIndex }) => {
  const countries = getCountriesOptions();
  const addressTypes = addressTypeOptions();
  const {
    watch,
    formState: { errors },
    setValue,
    control,
    register
  } = useFormContext();

  const getFirstAddressType = watch('entity.geographic_addresses[0].address_type');

  useEffect(() => {
    if (!getFirstAddressType) {
      setValue(`entity.geographic_addresses[0].address_type`, addressTypeEnum.ADDRESS_TYPE_BIZZ);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [getFirstAddressType]);

  return (
    <>
      <VStack spacing={3.5} align="start">
        <InputFormControl
          formHelperText={t`Address line 1 e.g. building name/number, street name (required)`}
          controlId={`${name}[${rowIndex}].address_line[0]`}
          isInvalid={!!getValueByPathname(errors, `${name}[${rowIndex}].address_line[0]`)}
          {...register(`${name}[${rowIndex}].address_line[0]`)}
          data-testid="address_line[0]"
        />

        <InputFormControl
          formHelperText={t`Address line 2 e.g. apartment or suite number`}
          controlId={`${name}[${rowIndex}].address_line[1]`}
          isInvalid={!!getValueByPathname(errors, `${name}[${rowIndex}].address_line[1]`)}
          {...register(`${name}[${rowIndex}].address_line[1]`)}
          data-testid="address_line[1]"
        />

        <InputFormControl
          formHelperText={t`City / Town / Municipality`}
          controlId="town_name"
          isInvalid={!!getValueByPathname(errors, `${name}[${rowIndex}].town_name`)}
          {...register(`${name}[${rowIndex}].town_name`)}
          data-testid="city"
        />

        <Grid templateColumns={{ base: '1fr', md: 'repeat(2, 1fr)' }} gap={6} width="100%">
          <GridItem>
            <InputFormControl
              formHelperText={t`Region / Province / State (required)`}
              controlId="country_sub_division"
              isInvalid={!!getValueByPathname(errors, `${name}[${rowIndex}].country_sub_division`)}
              {...register(`${name}[${rowIndex}].country_sub_division`)}
              data-testid="state"
            />
          </GridItem>
          <GridItem>
            <InputFormControl
              formHelperText={t`Postal Code / Postcode / ZIP Code`}
              controlId="post_code"
              isInvalid={!!getValueByPathname(errors, `${name}[${rowIndex}].post_code`)}
              {...register(`${name}[${rowIndex}].post_code`)}
              data-testid="post_code"
            />
          </GridItem>
        </Grid>

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
                  formHelperText={t`Country (required)`}
                  controlId="country"
                  data-testid="country"
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
                  defaultValue={addressTypeEnum.ADDRESS_TYPE_BIZZ}
                  isInvalid={!!getValueByPathname(errors, `${name}[${rowIndex}].address_type`)}
                  value={addressTypes.find((option) => option.value === field.value)}
                  onChange={(newValue: any) => field.onChange(newValue.value)}
                  options={addressTypes}
                  formHelperText={t`Address Type (required)`}
                  controlId="address_type"
                  data-testid="address_type"
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
