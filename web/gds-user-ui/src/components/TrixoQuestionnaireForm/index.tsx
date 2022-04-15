import { Grid, GridItem, Heading, Text, VStack } from '@chakra-ui/react';
import OtherJuridictions from 'components/OtherJuridictions';
import Regulations from 'components/Regulations';
import SwitchFormControl from 'components/SwitchFormControl';
import InputFormControl from 'components/ui/InputFormControl';
import SelectFormControl from 'components/ui/SelectFormControl';
import { getCountriesOptions } from 'constants/countries';
import { getCurrenciesOptions, getFinancialTransfertsPermittedOptions } from 'constants/trixo';
import FormLayout from 'layouts/FormLayout';
import { Controller, useFormContext } from 'react-hook-form';
import { useEffect } from 'react';

const TrixoQuestionnaireForm: React.FC = () => {
  const { register, control, getValues, setValue, watch } = useFormContext();
  const countries = getCountriesOptions();
  const financialTransfertsOptions = getFinancialTransfertsPermittedOptions();
  const currencies = getCurrenciesOptions();
  const getHasRequiredRegulatoryProgram = watch('trixo.has_required_regulatory_program');
  const getMustComplyRegulations = watch('trixo.must_comply_travel_rule');

  const getCountryFromLegalAddress = watch('entity.country_of_registration');

  useEffect(() => {
    if (getCountryFromLegalAddress) {
      setValue(`trixo.primary_national_jurisdiction`, getCountryFromLegalAddress);
    }
  }, [getCountryFromLegalAddress]);

  return (
    <FormLayout spacing={4}>
      <Controller
        control={control}
        name="trixo.primary_national_jurisdiction"
        render={({ field }) => (
          <SelectFormControl
            options={countries}
            controlId="primaryNationalJuridiction"
            label="Primary National Jurisdiction"
            ref={field.ref}
            name={field.name}
            value={countries.find((option) => option.value === field.value)}
            onChange={(newValue: any) => field.onChange(newValue.value)}
          />
        )}
      />
      <InputFormControl
        controlId="nameOfPrimaryRegulator"
        label="Name of Primary Regulator"
        formHelperText="The name of primary regulator or supervisory authority for your national jurisdiction"
        {...register('trixo.primary_regulator')}
      />
      <VStack align="start">
        <Heading size="md">Other Jurisdictions</Heading>
        <Text>Please add any other regulatory jurisdictions your organization complies with.</Text>

        <OtherJuridictions name={'trixo.other_jurisdictions'} />
      </VStack>
      <Controller
        control={control}
        name="trixo.financial_transfers_permitted"
        render={({ field }) => (
          <SelectFormControl
            ref={field.ref}
            name={field.name}
            options={financialTransfertsOptions}
            value={financialTransfertsOptions.find((option) => option.value === field.value)}
            onChange={(newValue: any) => field.onChange(newValue.value)}
            label="Is your organization permitted to send and/or receive transfers of virtual assets in the jurisdictions in which it operates?"
            controlId="financial_transfers_permitted"
          />
        )}
      />
      <VStack align="start">
        <Heading size="md">CDD & Travel Rule Policies</Heading>

        <Controller
          control={control}
          name="trixo.has_required_regulatory_program"
          render={({ field }) => (
            <SelectFormControl
              ref={field.ref}
              name={field.name}
              formatOptionLabel={(data: any) => {
                return data.value === 'partial' ? `${data.label} implemented` : data.label;
              }}
              options={financialTransfertsOptions}
              value={financialTransfertsOptions.find((option) => option.value === field.value)}
              onChange={(newValue: any) => field.onChange(newValue.value)}
              label=" Does your organization have a programme that sets minimum AML, CFT,
            KYC/CDD and Sanctions standards per the requirements of the
            jurisdiction(s) regulatory regimes where it is
            licensed/approved/registered?"
              controlId="has_required_regulatory_program"
            />
          )}
        />

        <VStack>
          <Text>
            Does your organization conduct KYC/CDD before permitting its customers to send/receive
            virtual asset transfers?
          </Text>
          <SwitchFormControl
            label="Conducts KYC before virtual asset transfers"
            controlId="conducts_customer_kyc"
            {...register('trixo.conducts_customer_kyc')}
          />
        </VStack>
      </VStack>
      {getHasRequiredRegulatoryProgram && getHasRequiredRegulatoryProgram === 'yes' && (
        <VStack align="start" w="100%">
          <Text>At what threshold and currency does your organization conduct KYC?</Text>
          <Grid templateColumns={{ base: '1fr 1fr', md: '2fr 1fr' }} gap={6} width="100%">
            <GridItem>
              <InputFormControl
                type="number"
                label=""
                controlId="kyc_threshold"
                {...register('trixo.kyc_threshold')}
              />
            </GridItem>
            <GridItem>
              <Controller
                control={control}
                name="trixo.kyc_threshold_currency"
                render={({ field }) => (
                  <SelectFormControl
                    ref={field.ref}
                    name={field.name}
                    options={currencies}
                    value={currencies.find((option) => option.value === field.value)}
                    onChange={(newValue: any) => field.onChange(newValue.value)}
                    controlId="trixo.kyc_threshold_currency"
                  />
                )}
              />
            </GridItem>
          </Grid>
          <Text fontSize="sm" color="whiteAlpha.600" mt="0 !important">
            Threshold to conduct KYC before permitting the customer to send/receive virtual asset
            transfers
          </Text>
        </VStack>
      )}
      ;
      <VStack>
        <Text>
          Is your organization required to comply with the application of the Travel Rule standards
          in the jurisdiction(s) where it is licensed/approved/registered?
        </Text>
        <SwitchFormControl
          label="Must comply with Travel Rule"
          controlId="trixo.must_comply_travel_rule"
          {...register('trixo.must_comply_travel_rule')}
        />
      </VStack>
      <VStack align="start" w="100%">
        <Heading size="md">Applicable Regulations</Heading>
        <Text fontSize="sm">
          Please specify the applicable regulation(s) for Travel Rule standards compliance, e.g.
          "FATF Recommendation 16"
        </Text>
        <Regulations register={register} name={`trixo.applicable_regulations`} control={control} />
      </VStack>
      {getMustComplyRegulations && (
        <VStack align="start" w="100%">
          <Text>What is the minimum threshold for Travel Rule compliance?</Text>
          <Grid templateColumns={{ base: '1fr 1fr', md: '2fr 1fr' }} gap={6} width="100%">
            <GridItem>
              <InputFormControl
                type="number"
                label=""
                controlId="compliance_threshold"
                {...register('trixo.compliance_threshold')}
              />
            </GridItem>
            <GridItem>
              <Controller
                control={control}
                name="trixo.compliance_threshold_currency"
                render={({ field }) => (
                  <SelectFormControl
                    ref={field.ref}
                    name={field.name}
                    options={currencies}
                    value={currencies.find((option) => option.value === field.value)}
                    onChange={(newValue: any) => field.onChange(newValue.value)}
                    controlId="trixo.compliance_threshold_currency"
                  />
                )}
              />
            </GridItem>
          </Grid>
          <Text fontSize="sm" mt="0 !important">
            The minimum threshold above which your organization is required to collect/send Travel
            Rule information.
          </Text>
        </VStack>
      )}
      <VStack align="start" w="100%" spacing={4}>
        <Heading size="md">Data Protection Policies</Heading>
        <VStack align="start" w="100%">
          <Text>Is your organization required by law to safeguard PII?</Text>
          <SwitchFormControl
            label="Must safeguard PII"
            controlId="must_safeguard_pii"
            {...register('trixo.must_safeguard_pii')}
          />
        </VStack>
        <VStack align="start" w="100%">
          <Text>
            Does your organization secure and protect PII, including PII received from other VASPs
            under the Travel Rule?
          </Text>
          <SwitchFormControl
            label="Safeguards PII"
            controlId="safeguards_pii"
            {...register('trixo.safeguards_pii')}
          />
        </VStack>
      </VStack>
    </FormLayout>
  );
};

export default TrixoQuestionnaireForm;
