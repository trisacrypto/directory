import { VStack, chakra } from '@chakra-ui/react';
import InputFormControl from 'components/ui/InputFormControl';
import SelectFormControl from 'components/ui/SelectFormControl';
import { getBusinessCategoryOptions, vaspCategories } from 'constants/basic-details';
import { Controller, FormProvider, useForm } from 'react-hook-form';
import { t } from '@lingui/macro';
import { useLanguageProvider } from 'contexts/LanguageContext';
import { useEffect } from 'react';
import FormLayout from 'layouts/FormLayout';
import formatDate from 'utils/formate-date';
import { yupResolver } from '@hookform/resolvers/yup';
import { basicDetailsValidationSchema } from 'modules/dashboard/certificate/lib/basicDetailsValidationSchema';
import StepButtons from 'components/StepsButtons';
import { isProdEnv } from 'application/config';
import { DevTool } from '@hookform/devtools';
import { useUpdateCertificateStep } from 'hooks/useUpdateCertificateStep';
import { useFetchCertificateStep } from 'hooks/useFetchCertificateStep';
import useCertificateStepper from 'hooks/useCertificateStepper';
import { StepEnum } from 'types/enums';
import { StepsIndexes } from 'constants/steps';
interface BasicDetailsFormProps {
  data?: any;
  isLoading?: boolean;
}
const BasicDetailsForm: React.FC<BasicDetailsFormProps> = ({ data }) => {
  const { certificateStep } = useFetchCertificateStep({
    key: StepEnum.BASIC
  });
  const { updateCertificateStep, updatedCertificateStep } = useUpdateCertificateStep();
  const resolver = yupResolver(basicDetailsValidationSchema);
  const options = getBusinessCategoryOptions();
  const { currentState, nextStep, updateIsDirty } = useCertificateStepper();

  const [language] = useLanguageProvider();

  useEffect(() => {}, [language]);

  const methods = useForm({
    defaultValues: certificateStep?.form || data,
    resolver
  });

  const {
    register,
    formState: { errors, isDirty },
    control
  } = methods;

  useEffect(() => {
    updateIsDirty(isDirty, StepsIndexes.BASIC_DETAILS);
  }, [isDirty, updateIsDirty]);

  const handleNextStepClick = () => {
    if (!isDirty) {
      nextStep(updatedCertificateStep?.errors ?? certificateStep?.errors);
    } else {
      const payload = {
        step: StepEnum.BASIC,
        form: {
          ...methods.getValues(),
          state: currentState()
        } as any
      };

      updateCertificateStep(payload);
      nextStep(updatedCertificateStep);
    }
  };

  return (
    <FormLayout>
      <FormProvider {...methods}>
        <chakra.form onSubmit={methods.handleSubmit(handleNextStepClick)}>
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
            <StepButtons handleNextStep={handleNextStepClick} isFirstStep={true} />
          </VStack>
        </chakra.form>
        {!isProdEnv ? <DevTool control={methods.control} /> : null}
      </FormProvider>
    </FormLayout>
  );
};

export default BasicDetailsForm;
