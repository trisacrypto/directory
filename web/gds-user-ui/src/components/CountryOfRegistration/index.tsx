import { Heading } from '@chakra-ui/react';
import { t } from '@lingui/macro';
import { Trans } from '@lingui/react';
import SelectFormControl from 'components/ui/SelectFormControl';
import { getCountriesOptions } from 'constants/countries';
import FormLayout from 'layouts/FormLayout';
import { Controller, useFormContext } from 'react-hook-form';

const CountryOfRegistration: React.FC = () => {
  const {
    control,
    formState: { errors }
  } = useFormContext();
  const countries = getCountriesOptions();

  return (
    <FormLayout>
      <Heading size="md">
        <Trans id="Country of Registration">Country of Registration</Trans>
      </Heading>
      <Controller
        control={control}
        name="entity.country_of_registration"
        render={({ field }) => (
          <SelectFormControl
            ref={field.ref}
            label=""
            placeholder={t`Select a country`}
            isInvalid={!!errors?.entity?.country_of_registration}
            formHelperText={errors?.entity?.country_of_registration?.message}
            controlId="entity.country_of_registration"
            options={countries}
            name={field.name}
            value={countries.find((option) => option.value === field.value)}
            onChange={(newValue: any) => field.onChange(newValue.value)}
          />
        )}
      />
    </FormLayout>
  );
};

export default CountryOfRegistration;
