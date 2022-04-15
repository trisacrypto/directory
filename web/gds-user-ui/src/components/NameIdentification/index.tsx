import React, { useState, useEffect, useRef } from 'react';
import { Heading, Link, Text } from '@chakra-ui/react';
import InputFormControl from 'components/ui/InputFormControl';
import SelectFormControl from 'components/ui/SelectFormControl';
import { getCountriesOptions } from 'constants/countries';
import { getNationalIdentificationOptions } from 'constants/national-identification';
import FormLayout from 'layouts/FormLayout';
import { Controller, useFormContext } from 'react-hook-form';
import { getRegistrationAuthoritiesOptions, getValueByPathname } from 'utils/utils';

interface NationalIdentificationProps {}

const NationalIdentification: React.FC<NationalIdentificationProps> = () => {
  const {
    register,
    control,
    watch,
    formState: { errors },
    setValue,
    clearErrors
  } = useFormContext();
  const nationalIdentificationOptions = getNationalIdentificationOptions();
  const countries = getCountriesOptions();
  // const registrationAuthority = getRegistrationAuthoritiesOptions();
  const NationalIdentificationType = watch(
    'entity.national_identification.national_identifier_type'
  );

  // eslint-disable-next-line prefer-const
  let inputRegRef = useRef<any>();

  useEffect(() => {
    if (NationalIdentificationType === 'NATIONAL_IDENTIFIER_TYPE_CODE_LEIX') {
      setValue('entity.national_identification.registration_authority', '');
      clearErrors('entity.national_identification.registration_authority');
      inputRegRef?.current?.clear();
    }
  }, [NationalIdentificationType]);
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
        isInvalid={!!errors?.entity?.national_identification?.national_identifier}
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
            isInvalid={
              !!getValueByPathname(
                errors,
                'entity.national_identification.national_identifier_type'
              )
            }
            formHelperText={
              getValueByPathname(errors, 'entity.national_identification.national_identifier_type')
                ? getValueByPathname(
                    errors,
                    'entity.national_identification.national_identifier_type'
                  ).message
                : null
            }
            value={nationalIdentificationOptions.find((option) => option.value === field.value)}
            onChange={(newValue: any) => field.onChange(newValue.value)}
            label="Identification Type"
            controlId="identification_type"
          />
        )}
      />

      {NationalIdentificationType !== 'NATIONAL_IDENTIFIER_TYPE_CODE_LEIX' && (
        <Controller
          control={control}
          name="entity.national_identification.country_of_issue"
          render={({ field }) => (
            <SelectFormControl
              ref={field.ref}
              options={countries}
              value={countries.find((option) => option.value === field.value)}
              onChange={(newValue: any) => field.onChange(newValue.value)}
              isInvalid={!!errors?.entity?.national_identification?.country_of_issue}
              label="Country of Issue"
              controlId="country_of_issue"
              formHelperText={errors?.entity?.national_identification?.country_of_issue?.message}
            />
          )}
        />
      )}

      {/* <Controller
        control={control}
        name="entity.national_identification.registration_authority"
        render={({ field }) => (
          <SelectFormControl
            ref={field.ref}
            options={registrationAuthority}
            value={registrationAuthority.find((option) => option.value === field.value)}
            onChange={(newValue: any) => field.onChange(newValue.value)}
            label="Registration Authority"
            controlId="registration_authority"
            isInvalid={!!errors?.entity?.national_identification?.registration_authority}
            isDisabled={NationalIdentificationType === 'NATIONAL_IDENTIFIER_TYPE_CODE_LEIX'}
            formHelperText={
              errors?.entity?.national_identification?.registration_authority?.message ||
              'If the identifier is an LEI number, enter the ID used in the GLEIF Registration Authorities List.'
            }
          />
        )}
      /> */}

      <InputFormControl
        label="Registration Authority"
        controlId="registration_authority"
        inputRef={inputRegRef && undefined}
        isInvalid={!!errors?.entity?.national_identification?.registration_authority}
        isRequired={NationalIdentificationType !== 'NATIONAL_IDENTIFIER_TYPE_CODE_LEIX' && false}
        isDisabled={NationalIdentificationType === 'NATIONAL_IDENTIFIER_TYPE_CODE_LEIX'}
        formHelperText="If the identifier is an LEI number, enter the ID used in the GLEIF Registration Authorities List."
        {...register('entity.national_identification.registration_authority')}
      />
    </FormLayout>
  );
};

export default NationalIdentification;
