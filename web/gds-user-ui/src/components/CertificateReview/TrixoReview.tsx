import React, { useEffect } from 'react';
import { useColorModeValue } from '@chakra-ui/react';
import { useSelector, RootStateOrAny } from 'react-redux';
import { TStep } from 'utils/localStorageHelper';
import useCertificateStepper from 'hooks/useCertificateStepper';
import { t } from '@lingui/macro';
import { getRegistrationDefaultValue } from 'modules/dashboard/registration/utils';
interface TrixoReviewProps {
  data?: any;
}
import TrixoReviewDataTable from './TrixoReviewDataTable';
import CertificateReviewHeader from './CertificateReviewHeader';
import CertificateReviewLayout from './CertificateReviewLayout';
interface TrixoReviewProps {}

const TrixoReview: React.FC<TrixoReviewProps> = (props) => {
  const { jumpToStep } = useCertificateStepper();
  const steps: TStep[] = useSelector((state: RootStateOrAny) => state.stepper.steps);
  const [trixo, setTrixo] = React.useState<any>({});
  const textColor = useColorModeValue('gray.800', '#F7F8FC');
  const getColorScheme = (status: string | boolean) => {
    if (status === 'yes' || status === true) {
      return 'green';
    } else {
      return 'orange';
    }
  };
  useEffect(() => {
    const fetchData = async () => {
      const getStepperData = await getRegistrationDefaultValue();
      const stepData = {
        ...getStepperData.trixo
      };
      setTrixo(stepData);
    };
    fetchData();
  }, [steps]);

  return (
    <CertificateReviewLayout>
      <CertificateReviewHeader title={t`Section 5: TRIXO Questionnaire`} step={5} />
      <TrixoReviewDataTable data={trixo} />
    </CertificateReviewLayout>
  );
};

export default TrixoReview;
