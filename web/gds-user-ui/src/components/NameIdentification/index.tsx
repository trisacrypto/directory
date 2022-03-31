import { Heading, Link, Text } from '@chakra-ui/react';
import InputFormControl from 'components/ui/InputFormControl';
import SelectFormControl from 'components/ui/SelectFormControl';
import { getCountriesOptions } from 'constants/countries';
import { getNationalIdentificationOptions } from 'constants/national-identification';
import FormLayout from 'layouts/FormLayout';
import { Controller, useFormContext } from 'react-hook-form';

type NationalIdentificationProps = {};

const NationalIdentification: React.FC<NationalIdentificationProps> = () => {
  const { register, control } = useFormContext();
  const nationalIdentificationOptions = getNationalIdentificationOptions();
  const countries = getCountriesOptions();

  return (
    <FormLayout>
      <Heading size="md">National Identification</Heading>
      <Text>
        Please supply a valid national identification number. TRISA recommends the use of LEI
        numbers. For more information, please visit{' '}
        <Link href="https://gleif.org" color="blue.500" isExternal>
          GLEIF.org
        </Link>
      </Text>
      <InputFormControl
        label="Identification Number"
        controlId="identification_number"
        formHelperText="An identifier issued by an appropriate issuing authority"
        {...register('entity.national_identification.national_identifier')}
      />
      <Controller
        control={control}
        name={'entity.national_identification.national_identifier_type'}
        render={({ field }) => (
          <SelectFormControl
            ref={field.ref}
            options={nationalIdentificationOptions}
            value={nationalIdentificationOptions.find((option) => option.value === field.value)}
            onChange={(newValue: any) => field.onChange(newValue.value)}
            label="Identification Type"
            controlId="identification_type"
          />
        )}
      />

      <Controller
        control={control}
        name="entity.national_identification.country_of_issue"
        render={({ field }) => (
          <SelectFormControl
            ref={field.ref}
            options={countries}
            value={countries.find((option) => option.value === field.value)}
            onChange={(newValue: any) => field.onChange(newValue.value)}
            label="Country of Issue"
            controlId="country_of_issue"
          />
        )}
      />
      <InputFormControl
        label="Registration Authority"
        controlId="registration_authority"
        formHelperText="If the identifier is an LEI number, enter the ID used in the GLEIF Registration Authorities List."
        {...register('entity.national_identification.registration_authority')}
      />
    </FormLayout>
  );
};

export default NationalIdentification;
