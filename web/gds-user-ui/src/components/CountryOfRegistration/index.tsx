import { Heading } from '@chakra-ui/react';
import InputFormControl from 'components/ui/InputFormControl';
import FormLayout from 'layouts/FormLayout';
import { useFormContext } from 'react-hook-form';

type CountryOfRegistrationProps = {};
const CountryOfRegistration: React.FC<CountryOfRegistrationProps> = () => {
  const { register } = useFormContext();
  return (
    <FormLayout>
      <Heading size="md">Country of Registration</Heading>
      <InputFormControl
        controlId="entity.country_of_registration"
        {...register('entity.country_of_registration')}
      />
    </FormLayout>
  );
};

export default CountryOfRegistration;
