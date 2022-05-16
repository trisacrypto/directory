import React, { useEffect, useRef } from 'react';
import { Heading, Link, Text } from '@chakra-ui/react';
import InputFormControl from 'components/ui/InputFormControl';
import SelectFormControl from 'components/ui/SelectFormControl';
import { getCountriesOptions } from 'constants/countries';
import {
  getNationalIdentificationOptions,
  disabledIdentifiers
} from 'constants/national-identification';
import FormLayout from 'layouts/FormLayout';
import { Controller, useFormContext } from 'react-hook-form';
import {
  getRegistrationAuthorities,
  getRegistrationAuthoritiesOptions,
  getValueByPathname
} from 'utils/utils';

interface NationalIdentificationProps {}

const RegistrationAuthorityFormHelperText = ({ option }: { option: string }) => {
  const registrationAuthority = React.useMemo(() => {
    const _authority = getRegistrationAuthorities().find(
      (authority) => authority.option === option
    );
    return _authority;
  }, [option]);

  return (
    <Text>
      For identifiers other than LEI specify the registration authority from the following list. See{' '}
      <Link
        href="https://www.gleif.org/en/about-lei/code-lists/gleif-registration-authorities-list"
        color="blue.500"
        isExternal>
        GLEIF Registration Authorities
      </Link>{' '}
      for more details on how to look up a registration authority. If in doubt, use RA777777 -
      "General Government Entities" which specifies the default registration authority for your
      country of registration.
      {registrationAuthority?.website && (
        <Text color={'#1a202c'} fontSize="sm" mt={3}>
          Website:{' '}
          <Link color={'blue.400'} isExternal>
            {registrationAuthority?.website}
          </Link>
        </Text>
      )}
      {registrationAuthority?.comments && (
        <Text color={'#1a202c'} fontSize="sm" mt={2}>
          Comments:{' '}
          <Text fontStyle={'italic'} as="span">
            {registrationAuthority?.comments}
          </Text>
        </Text>
      )}
    </Text>
  );
};

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

  const NationalIdentificationType = watch(
    'entity.national_identification.national_identifier_type'
  );
  const getCountryOfRegistration = watch('entity.country_of_registration');

  const registrationAuthority = getRegistrationAuthoritiesOptions(getCountryOfRegistration);
  const getRegistrationAuthority = () => {
    // setValue('entity.national_identification.registration_authority', 'RA777777');
    return getRegistrationAuthoritiesOptions(getCountryOfRegistration);
  };
  // eslint-disable-next-line prefer-const
  let inputRegRef = useRef<any>();

  useEffect(() => {
    if (NationalIdentificationType === 'NATIONAL_IDENTIFIER_TYPE_CODE_LEIX') {
      setValue('entity.national_identification.registration_authority', 'RA777777');
      clearErrors('entity.national_identification.registration_authority');

      inputRegRef?.current?.clear();
    }
    if (
      NationalIdentificationType !== 'NATIONAL_IDENTIFIER_TYPE_CODE_LEIX' &&
      !disabledIdentifiers.includes(NationalIdentificationType)
    ) {
      setValue('entity.national_identification.registration_authority', 'RA777777');
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
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
      {disabledIdentifiers.includes(NationalIdentificationType) && (
        <>
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
                isDisabled={!disabledIdentifiers.includes(NationalIdentificationType)}
                label="Country of Issue"
                controlId="country_of_issue"
                formHelperText={
                  errors?.entity?.national_identification?.country_of_issue?.message ||
                  'Country of Issue is reserved for National Identifiers of Natural Persons'
                }
              />
            )}
          />
        </>
      )}
      {NationalIdentificationType !== 'NATIONAL_IDENTIFIER_TYPE_CODE_LEIX' && (
        <>
          <Controller
            control={control}
            name="entity.national_identification.registration_authority"
            render={({ field }) => (
              <SelectFormControl
                ref={field.ref}
                options={getRegistrationAuthority()}
                value={registrationAuthority.find((option) => option.value === field.value)}
                onChange={(newValue: any) => field.onChange(newValue.value)}
                label="Registration Authority"
                controlId="registration_authority"
                isInvalid={!!errors?.entity?.national_identification?.registration_authority}
                isDisabled={NationalIdentificationType === 'NATIONAL_IDENTIFIER_TYPE_CODE_LEIX'}
                formHelperText={
                  errors?.entity?.national_identification?.registration_authority?.message || (
                    <RegistrationAuthorityFormHelperText option={field.value} />
                  )
                }
              />
            )}
          />
        </>
      )}
    </FormLayout>
  );
};

export default NationalIdentification;
