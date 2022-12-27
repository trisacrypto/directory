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
import { t } from '@lingui/macro';
import { Trans } from '@lingui/react';

const TrixoQuestionnaireForm: React.FC = () => {
  const { register, control, setValue, watch } = useFormContext();
  const countries = getCountriesOptions();
  const financialTransfertsOptions = getFinancialTransfertsPermittedOptions();
  const currencies = getCurrenciesOptions();
  const getHasRequiredRegulatoryProgram = watch('trixo.has_required_regulatory_program');
  const getMustComplyRegulations = watch('trixo.must_comply_travel_rule');
  const getApplicableRegulations = watch('trixo.applicable_regulations');
  const getCountryFromLegalAddress = watch('entity.country_of_registration');
  const getComplianceThreshold = watch('trixo.compliance_threshold');
  const getKycThreshold = watch('trixo.kyc_threshold');

  useEffect(() => {
    if (getCountryFromLegalAddress) {
      setValue(`trixo.primary_national_jurisdiction`, getCountryFromLegalAddress);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [getCountryFromLegalAddress]);

  // set default value if getMustComplyRegulations or  is false
  useEffect(() => {
    if (!getMustComplyRegulations) {
      setValue(`trixo.compliance_threshold`, 0);
    }
    if (
      !getHasRequiredRegulatoryProgram ||
      getHasRequiredRegulatoryProgram === 'no' ||
      getHasRequiredRegulatoryProgram === 'partial'
    ) {
      setValue(`trixo.kyc_threshold`, 0);
    }
    // if applicable regulations is empty, set default value
    if (getApplicableRegulations?.length === 0) {
      setValue(`trixo.applicable_regulations`, ['FATF Recommendation 16']);
    }
  }, [
    getMustComplyRegulations,
    getHasRequiredRegulatoryProgram,
    setValue,
    getApplicableRegulations
  ]);

  useEffect(() => {
    const regExp = /^0[0-9].*$/;
    if (getComplianceThreshold !== 0) {
      if (regExp.test(getComplianceThreshold)) {
        setValue(`trixo.compliance_threshold`, getComplianceThreshold.replace(/^0+/, ''));
      }
    }
    if (getKycThreshold !== 0) {
      if (regExp.test(getKycThreshold)) {
        setValue(`trixo.kyc_threshold`, getKycThreshold.replace(/^0+/, ''));
      }
    }
  }, [getKycThreshold, getComplianceThreshold, setValue]);

  return (
    <FormLayout spacing={5}>
      <Controller
        control={control}
        name="trixo.primary_national_jurisdiction"
        render={({ field }) => (
          <SelectFormControl
            options={countries}
            controlId="primaryNationalJuridiction"
            label={t`Primary National Jurisdiction`}
            ref={field.ref}
            name={field.name}
            value={countries.find((option) => option.value === field.value)}
            onChange={(newValue: any) => field.onChange(newValue.value)}
          />
        )}
      />
      <InputFormControl
        controlId="nameOfPrimaryRegulator"
        label={t`Name of Primary Regulator`}
        formHelperText={t`The name of primary regulator or supervisory authority for your national jurisdiction`}
        {...register('trixo.primary_regulator')}
      />
      <VStack align="start" w="100%">
        <Heading size="md">
          <Trans id="Other Jurisdictions">Other Jurisdictions</Trans>
        </Heading>
        <Text>
          <Trans id="Please add any other regulatory jurisdictions your organization complies with.">
            Please add any other regulatory jurisdictions your organization complies with.
          </Trans>
        </Text>

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
            label={t`Is your organization permitted to send and/or receive transfers of virtual assets in the jurisdictions in which it operates?`}
            controlId="financial_transfers_permitted"
          />
        )}
      />
      <VStack align="start">
        <Heading size="md">
          <Trans id="CDD & Travel Rule Policies">CDD & Travel Rule Policies</Trans>
        </Heading>

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
              label={t`Does your organization have a programme that sets minimum Anti-Money
              Laundering (AML), Countering the Financing of Terrorism (CFT), Know your
              Counterparty/Customer Due Diligence (KYC/CDD) and Sanctions standards per the requirements of the jurisdiction(s) regulatory regimes where it is licensed/approved/registered?`}
              controlId="has_required_regulatory_program"
            />
          )}
        />

        <VStack>
          <Text>
            <Trans id="Does your organization conduct KYC/CDD before permitting its customers to send/receive virtual asset transfers?">
              Does your organization conduct KYC/CDD before permitting its customers to send/receive
              virtual asset transfers?
            </Trans>
          </Text>
          <SwitchFormControl
            label={t`Conducts KYC before virtual asset transfers`}
            controlId="conducts_customer_kyc"
            {...register('trixo.conducts_customer_kyc')}
          />
        </VStack>
      </VStack>
      {getHasRequiredRegulatoryProgram && getHasRequiredRegulatoryProgram === 'yes' && (
        <VStack align="start" w="100%">
          <Text>
            <Trans id="At what threshold and currency does your organization conduct KYC checks?">
              At what threshold and currency does your organization conduct KYC checks?
            </Trans>
          </Text>
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
          <Text fontSize="sm" color="gray.500" mt="0 !important">
            <Trans id="Threshold to conduct KYC before permitting the customer to send/receive virtual asset transfers">
              Threshold to conduct KYC before permitting the customer to send/receive virtual asset
              transfers
            </Trans>
          </Text>
        </VStack>
      )}
      ;
      <VStack>
        <Text>
          <Trans id="Is your organization required to comply with the application of the Travel Rule standards in the jurisdiction(s) where it is licensed/approved/registered?">
            Is your organization required to comply with the application of the Travel Rule
            standards in the jurisdiction(s) where it is licensed/approved/registered?
          </Trans>
        </Text>
        <SwitchFormControl
          label={t`Must comply with Travel Rule`}
          controlId="trixo.must_comply_travel_rule"
          {...register('trixo.must_comply_travel_rule')}
        />
      </VStack>
      {getMustComplyRegulations && (
        <VStack align="start" w="100%">
          <Text>
            <Trans id="What is the minimum threshold for Travel Rule compliance?">
              What is the minimum threshold for Travel Rule compliance?
            </Trans>
          </Text>
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
            <Trans id="The minimum threshold above which your organization is required to collect/send Travel Rule information.">
              The minimum threshold above which your organization is required to collect/send Travel
              Rule information.
            </Trans>
          </Text>
        </VStack>
      )}
      <VStack align="start" w="100%">
        <Heading size="md">
          <Trans id="Applicable Regulations">Applicable Regulations</Trans>
        </Heading>
        <Text fontSize="sm">
          <Trans id='Please specify the applicable regulation(s) for Travel Rule standards compliance, e.g."FATF Recommendation 16"'>
            Please specify the applicable regulation(s) for Travel Rule standards compliance, e.g.
            "FATF Recommendation 16"
          </Trans>
        </Text>
        <Regulations name={`trixo.applicable_regulations`} />
      </VStack>
      <VStack align="start" w="100%" spacing={4}>
        <Heading size="md">
          <Trans id="Data Protection Policies">Data Protection Policies</Trans>
        </Heading>
        <VStack align="start" w="100%">
          <Text>
            <Trans id="Is your organization required by law to safeguard Personally Identifiable Information (PII)?">
              Is your organization required by law to safeguard Personally Identifiable Information
              (PII)?
            </Trans>
          </Text>
          <SwitchFormControl
            label={t`Must safeguard PII`}
            controlId="must_safeguard_pii"
            {...register('trixo.must_safeguard_pii')}
          />
        </VStack>
        <VStack align="start" w="100%">
          <Text>
            <Trans id="Does your organization secure and protect PII, including PII received from other VASPs under the Travel Rule?">
              Does your organization secure and protect PII, including PII received from other VASPs
              under the Travel Rule?
            </Trans>
          </Text>
          <SwitchFormControl
            label={t`Safeguards PII`}
            controlId="safeguards_pii"
            {...register('trixo.safeguards_pii')}
          />
        </VStack>
      </VStack>
    </FormLayout>
  );
};

export default TrixoQuestionnaireForm;
