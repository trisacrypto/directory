import { Heading } from '@chakra-ui/react';
import SelectFormControl from 'components/ui/SelectFormControl';
import { getCountriesOptions } from 'constants/countries';
import FormLayout from 'layouts/FormLayout';
import { Controller, useFormContext } from 'react-hook-form';

type CountryOfRegistrationProps = {};
const CountryOfRegistration: React.FC<CountryOfRegistrationProps> = () => {
  const {
    control,
    formState: { errors }
  } = useFormContext();
  const countries = getCountriesOptions();

  return (
    <FormLayout>
      <Heading size="md">Country of Registration</Heading>
      <Controller
        control={control}
        name="entity.country_of_registration"
        render={({ field }) => (
          <SelectFormControl
            ref={field.ref}
            label=""
            placeholder="Select a country"
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
