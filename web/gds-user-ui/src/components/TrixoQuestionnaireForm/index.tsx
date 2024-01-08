import { useEffect, useState, useRef, Dispatch, SetStateAction } from 'react';
import { Box, Grid, GridItem, Heading, Text, VStack, chakra, useDisclosure } from '@chakra-ui/react';
import OtherJuridictions from 'components/OtherJuridictions';
import Regulations from 'components/Regulations';
import SwitchFormControl from 'components/SwitchFormControl';
import InputFormControl from 'components/ui/InputFormControl';
import SelectFormControl from 'components/ui/SelectFormControl';
import { getCountriesOptions } from 'constants/countries';
import { getCurrenciesOptions, getFinancialTransfertsPermittedOptions } from 'constants/trixo';
import FormLayout from 'layouts/FormLayout';
import { Controller, useForm, FormProvider } from 'react-hook-form';

import { t } from '@lingui/macro';
import { Trans } from '@lingui/react';
import { useUpdateCertificateStep } from 'hooks/useUpdateCertificateStep';

import StepButtons from 'components/StepsButtons';
import { StepEnum } from 'types/enums';
import { trixoQuestionnaireValidationSchema } from 'modules/dashboard/certificate/lib/trixoQuestionnaireValidationSchema';
import { yupResolver } from '@hookform/resolvers/yup';
import useCertificateStepper from 'hooks/useCertificateStepper';
import { StepsIndexes } from 'constants/steps';
import { isProdEnv } from 'application/config';
import { DevTool } from '@hookform/devtools';

interface TrixoFormProps {
  data: any;
  isLoading?: boolean;
  shouldResetForm?: boolean;
  onResetFormState?: Dispatch<SetStateAction<boolean>>;
}
const TrixoQuestionnaireForm: React.FC<TrixoFormProps> = ({
  data,
  shouldResetForm,
  onResetFormState
}) => {
  const { isOpen, onClose, onOpen } = useDisclosure();
  const [shouldShowResetFormModal, setShouldShowResetFormModal] = useState<boolean>(false);
  const { previousStep, nextStep, currentState, updateIsDirty } = useCertificateStepper();

  const resolver = yupResolver(trixoQuestionnaireValidationSchema);
  const methods = useForm({
    defaultValues: data,
    resolver,
    mode: 'onChange'
  });

  const {
    register,
    control,
    setValue,
    watch,
    formState: { isDirty },
    reset: resetForm
  } = methods;

  const previousStepRef = useRef<any>(false);
  const nextStepRef = useRef<any>(false);
  const countries = getCountriesOptions();
  const financialTransfertsOptions = getFinancialTransfertsPermittedOptions();
  const currencies = getCurrenciesOptions();
  const getHasRequiredRegulatoryProgram = watch('trixo.has_required_regulatory_program');
  const getMustComplyRegulations = watch('trixo.must_comply_travel_rule');
  const getApplicableRegulations = watch('trixo.applicable_regulations');
  const getCountryFromLegalAddress = watch('entity.country_of_registration');
  const getComplianceThreshold = watch('trixo.compliance_threshold');
  const getKycThreshold = watch('trixo.kyc_threshold');
  const getMustComplyRegulationsFromData = data?.trixo?.must_comply_travel_rule;
  const getHasRequiredRegulatoryProgramFromData = data?.trixo?.has_required_regulatory_program;

  const {
    updateCertificateStep,
    updatedCertificateStep,
    isUpdatingCertificateStep,
    wasCertificateStepUpdated,
    reset: resetMutation
  } = useUpdateCertificateStep();
  const onCloseModalHandler = () => {
    setShouldShowResetFormModal(false);
    onClose();
  };

  if (wasCertificateStepUpdated && nextStepRef.current) {
    resetMutation();
    // reset the form with the new values
    resetForm(updatedCertificateStep?.form, {
      keepValues: false
    });
    nextStep(updatedCertificateStep);
    nextStepRef.current = false;
  }

  if (wasCertificateStepUpdated && previousStepRef.current && !isUpdatingCertificateStep) {
    resetMutation();
    // reset the form with the new values
    resetForm(updatedCertificateStep?.form, {
      keepValues: false
    });
    previousStepRef.current = false;
    previousStep(updatedCertificateStep);
  }

  const handleNextStepClick = () => {
    if (
      isDirty ||
      getMustComplyRegulationsFromData !== getMustComplyRegulations ||
      getHasRequiredRegulatoryProgramFromData !== getHasRequiredRegulatoryProgram
    ) {
      const payload = {
        step: StepEnum.TRIXO,
        form: {
          ...methods.getValues(),
          state: currentState()
        } as any
      };

      updateCertificateStep(payload);
      nextStepRef.current = true;
    } else {
      nextStep({
        step: StepEnum.TRIXO,
        form: {
          ...methods.getValues(),
          state: currentState()
        } as any
      });
    }
  };

  const handlePreviousStepClick = () => {
    // isDirty is not working for the checkbox so we need to compare the values
    if (
      isDirty ||
      getMustComplyRegulationsFromData !== getMustComplyRegulations ||
      getHasRequiredRegulatoryProgramFromData !== getHasRequiredRegulatoryProgram
    ) {
      const payload = {
        step: StepEnum.TRIXO,
        form: {
          ...methods.getValues(),
          state: currentState()
        } as any
      };
      console.log('[] isDirty  payload Trixo', payload);

      updateCertificateStep(payload);
      previousStepRef.current = true;
    }
    previousStep(data);
  };

  const handleResetForm = () => {
    setShouldShowResetFormModal(true); // this will show the modal
  };

  const handleResetClick = () => {
    setShouldShowResetFormModal(false); // this will close the modal
  };

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

  // set default value if getMustComplyRegulations or  is false
  useEffect(() => {
    if (!getMustComplyRegulations) {
      setValue(`trixo.compliance_threshold`, 0);
    }
    if (
      !getHasRequiredRegulatoryProgram ||
      getHasRequiredRegulatoryProgram === 'no' ||
      getHasRequiredRegulatoryProgram === 'partially'
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
    if (shouldShowResetFormModal) {
      onOpen();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [shouldShowResetFormModal]);

  useEffect(() => {
    updateIsDirty(isDirty, StepsIndexes.TRIXO_QUESTIONNAIRE);
  }, [isDirty, updateIsDirty]);

  useEffect(() => {
    if (getCountryFromLegalAddress) {
      setValue(`trixo.primary_national_jurisdiction`, getCountryFromLegalAddress);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [getCountryFromLegalAddress]);

  // reset the form from the parent component
  useEffect(() => {
    if (shouldResetForm && onResetFormState) {
      resetForm(data);
      onResetFormState(false);
      window.location.reload();
    }
  }, [shouldResetForm, resetForm, data, onResetFormState]);

  return (
    <FormLayout spacing={5}>
      <FormProvider {...methods}>
        <chakra.form onSubmit={methods.handleSubmit(handleNextStepClick)} data-testid="trixo-form">
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
          <VStack align="start" w="100%" py={5}>
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

          <VStack data-testid="is-required-financial-transfers">
            <Controller
              control={control}
              name="trixo.financial_transfers_permitted"
              render={({ field }) => (
                <SelectFormControl
                  ref={field.ref}
                  name={field.name}
                  data-testid="financial_transfers_permitted"
                  options={financialTransfertsOptions}
                  value={financialTransfertsOptions.find((option) => option.value === field.value)}
                  onChange={(newValue: any) => field.onChange(newValue.value)}
                  label={t`Is your organization permitted to send and/or receive transfers of virtual assets in the jurisdictions in which it operates?`}
                  controlId="financial_transfers_permitted"
                />
              )}
            />
          </VStack>
          <VStack align="start" pt={5} data-testid="trixo-rule-policies">
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
                  options={financialTransfertsOptions.filter(
                    (option) => option.value !== 'partially'
                  )}
                  value={financialTransfertsOptions.find((option) => option.value === field.value)}
                  onChange={(newValue: any) => field.onChange(newValue.value)}
                  label={t`Does your organization have a programme that sets minimum Anti-Money
              Laundering (AML), Countering the Financing of Terrorism (CFT), Know your
              Counterparty/Customer Due Diligence (KYC/CDD) and Sanctions standards per the requirements of the jurisdiction(s) regulatory regimes where it is licensed/approved/registered?`}
                  controlId="has_required_regulatory_program"
                />
              )}
            />

            <VStack data-testid="trixo-kyc-before-virtual-asset-transfers">
              <Text>
                <Trans id="Does your organization conduct KYC/CDD before permitting its customers to send/receive virtual asset transfers?">
                  Does your organization conduct KYC/CDD before permitting its customers to
                  send/receive virtual asset transfers?
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
                  Threshold to conduct KYC before permitting the customer to send/receive virtual
                  asset transfers
                </Trans>
              </Text>
            </VStack>
          )}

          <VStack align="start" data-testid="trixo-must-comply-travel-rule">
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
            <VStack align="start" w="100%" data-testid="tx-minimum-threshold">
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
                  The minimum threshold above which your organization is required to collect/send
                  Travel Rule information.
                </Trans>
              </Text>
            </VStack>
          )}
          <VStack align="start" w="100%" pt={5}>
            <Heading size="md">
              <Trans id="Applicable Regulations">Applicable Regulations</Trans>
            </Heading>
            <Text fontSize="sm">
              <Trans id='Please specify the applicable regulation(s) for Travel Rule standards compliance, e.g."FATF Recommendation 16"'>
                Please specify the applicable regulation(s) for Travel Rule standards compliance,
                e.g. "FATF Recommendation 16"
              </Trans>
            </Text>
            <Regulations name={`trixo.applicable_regulations`} />
          </VStack>
          <VStack align="start" w="100%" spacing={4} pt={5}>
            <Heading size="md">
              <Trans id="Data Protection Policies">Data Protection Policies</Trans>
            </Heading>
            <VStack align="start" w="100%">
              <Text>
                <Trans id="Is your organization required by law to safeguard Personally Identifiable Information (PII)?">
                  Is your organization required by law to safeguard Personally Identifiable
                  Information (PII)?
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
                  Does your organization secure and protect PII, including PII received from other
                  VASPs under the Travel Rule?
                </Trans>
              </Text>
              <SwitchFormControl
                label={t`Safeguards PII`}
                controlId="safeguards_pii"
                {...register('trixo.safeguards_pii')}
              />
            </VStack>
          </VStack>
          <Box pt={5}>
            <StepButtons
            handleNextStep={handleNextStepClick}
            handlePreviousStep={handlePreviousStepClick}
            onResetModalClose={handleResetClick}
            isOpened={isOpen}
            handleResetForm={handleResetForm}
            resetFormType={StepEnum.TRIXO}
            onClosed={onCloseModalHandler}
            handleResetClick={handleResetClick}
            shouldShowResetFormModal={shouldShowResetFormModal}
            />
          </Box>
        </chakra.form>
        {!isProdEnv ? <DevTool control={methods.control} /> : null}
      </FormProvider>
    </FormLayout>
  );
};

export default TrixoQuestionnaireForm;
