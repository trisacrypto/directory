import { VStack } from '@chakra-ui/react';
import InputFormControl from 'components/ui/InputFormControl';
import SelectFormControl from 'components/ui/SelectFormControl';
import { getBusinessCategoryOptions, vaspCategories } from 'constants/basic-details';
import { Controller, useFormContext } from 'react-hook-form';
import { t } from '@lingui/macro';
import { useLanguageProvider } from 'contexts/LanguageContext';
import { useEffect } from 'react';
import FormLayout from 'layouts/FormLayout';
import formatDate from 'utils/formate-date';

const BasicDetailsForm: React.FC = () => {
  const options = getBusinessCategoryOptions();
  const {
    register,
    control,
    formState: { errors }
  } = useFormContext();
  const [language] = useLanguageProvider();

  useEffect(() => {}, [language]);

  return (
    <FormLayout>
      <VStack spacing={4} w="100%">
        <InputFormControl
          controlId="organization_name"
          data-testid="organization_name"
          label={t`Organization Name`}
          error="true"
          formHelperText={errors.organization_name?.message as string}
          isInvalid={!!errors.organization_name}
          inputProps={{ placeholder: 'VASP HOLDING LLC' }}
          isRequiredField
          {...register('organization_name')}
        />

        <InputFormControl
          controlId="website"
          data-testid="website"
          label={t`Website`}
          error="true"
          type="url"
          formHelperText={errors.website?.message as string}
          isInvalid={!!errors.website}
          inputProps={{ placeholder: 'https://example.com' }}
          isRequiredField
          {...register('website')}
        />

        <InputFormControl
          controlId="established_on"
          data-testid="established_on"
          label={t`Date of Incorporation / Establishment`}
          formHelperText={errors.established_on?.message as string}
          isInvalid={!!errors.established_on}
          inputProps={{
            placeholder: '21/01/2021',
            type: 'date',
            pattern: 'd{4}-d{2}-d{2}',
            min: '1800-01-01',
            max: formatDate()
          }}
          isRequiredField
          {...register('established_on')}
        />

        <Controller
          control={control}
          name="business_category"
          render={({ field }) => (
            <SelectFormControl
              data-testid="business_category"
              ref={field.ref}
              label={t`Business Category`}
              placeholder={t`Select business category`}
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
              label={t`VASP Category`}
              placeholder={t`Select VASP category`}
              controlId="vasp_categories"
              data-testid="vasp_categories"
              isMulti
              name={name}
              options={vaspCategories}
              onChange={(val: any) => onChange(val.map((c: any) => c.value))}
              value={value && vaspCategories.filter((c) => value.includes(c.value))}
              formHelperText={t`Please select as many categories needed to represent the types of virtual asset services your organization provides.`}
            />
          )}
        />
      </VStack>
    </FormLayout>
  );
};

export default BasicDetailsForm;
