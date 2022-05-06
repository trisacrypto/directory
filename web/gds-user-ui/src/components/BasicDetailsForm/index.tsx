import { VStack, Text, FormErrorMessage } from '@chakra-ui/react';
import InputFormControl from 'components/ui/InputFormControl';
import SelectFormControl from 'components/ui/SelectFormControl';
import { getBusinessCategoryOptions, vaspCategories } from 'constants/basic-details';
import { Controller, useFormContext } from 'react-hook-form';
import { Control, UseFormRegister } from 'react-hook-form/dist/types/form';
import { yupResolver } from '@hookform/resolvers/yup';
import { ValidationSchema, getDefaultValue } from './validation';
import { useEffect } from 'react';
type BasicDetailsFormProps = {};

const BasicDetailsForm: React.FC<BasicDetailsFormProps> = () => {
  const options = getBusinessCategoryOptions();
  const {
    register,
    control,
    formState: { errors },
    getValues,
    watch,
    setValue
  } = useFormContext();

  // const getFirstLegalName = getValues('entity.name.name_identifiers')[0]?.legal_person_name;

  // const setLegalName = () => {
  //   if (getFirstLegalName && getFirstLegalName.length > 0) {
  //     setValue('organization', getFirstLegalName);
  //   }
  // };

  // useEffect(() => {
  //   setLegalName();
  // }, [getFirstLegalName]);
  return (
    <>
      <VStack spacing={4} w="100%">
        <InputFormControl
          controlId="organization_name"
          data-testid="organization_name"
          label="Organization Name"
          error="true"
          formHelperText={errors.organization_name?.message}
          isInvalid={!!errors.organization_name}
          inputProps={{ placeholder: 'VASP HOLDING LLC' }}
          {...register('organization_name')}
        />

        <InputFormControl
          controlId="website"
          data-testid="website"
          label="Website"
          error="true"
          type="url"
          formHelperText={errors.website?.message}
          isInvalid={!!errors.website}
          inputProps={{ placeholder: 'https://example.com' }}
          {...register('website')}
        />

        <InputFormControl
          controlId="established_on"
          data-testid="established_on"
          label="Date of Incorporation / Establishment"
          formHelperText={errors.established_on?.message}
          isInvalid={!!errors.established_on}
          inputProps={{ placeholder: '21/01/2021', type: 'date' }}
          {...register('established_on')}
        />

        <Controller
          control={control}
          name="business_category"
          render={({ field }) => (
            <SelectFormControl
              data-testid="business_category"
              ref={field.ref}
              label="Business Category"
              placeholder="Select business category"
              controlId="business_category"
              options={getBusinessCategoryOptions()}
              name={field.name}
              value={options.find((option) => option.value === field.value)}
              onChange={(newValue: any) => field.onChange(newValue.value)}
            />
          )}
        />

        <Controller
          control={control}
          name="vasp_categories"
          render={({ field: { value, onChange, name } }) => (
            <SelectFormControl
              label="VASP Category"
              placeholder="Select VASP category"
              controlId="vasp_categories"
              data-testid="vasp_categories"
              isMulti
              name={name}
              options={vaspCategories}
              onChange={(val: any) => onChange(val.map((c: any) => c.value))}
              value={value && vaspCategories.filter((c) => value.includes(c.value))}
              formHelperText="Please select as many categories needed to represent the types of virtual asset services your organization provides."
            />
          )}
        />
      </VStack>
    </>
  );
};

export default BasicDetailsForm;
