import { useSelector } from 'react-redux';

// NOTE: need some clean up.
import LegalPersonReviewDataTable from './LegalPersonReviewDataTable';
import CertificateReviewHeader from './CertificateReviewHeader';
import CertificateReviewLayout from './CertificateReviewLayout';
import useGetStepStatusByKey from './useGetStepStatusByKey';
import RequiredElementMissing from 'components/ErrorComponent/RequiredElementMissing';
import { t } from '@lingui/macro';
import { getCurrentState } from 'application/store/selectors/stepper';
const LegalPersonReview = () => {
  const currentStateValue = useSelector(getCurrentState);
  const { hasErrorField } = useGetStepStatusByKey(2);
  const legalPerson = {
    ...currentStateValue.data.entity
  };
  return (
    <CertificateReviewLayout>
      <CertificateReviewHeader step={2} title={t`Section 2: Legal Person`} />
      {hasErrorField ? <RequiredElementMissing elementKey={2} /> : false}
      <LegalPersonReviewDataTable data={legalPerson} />
    </CertificateReviewLayout>
  );
};
LegalPersonReview.defaultProps = {
  data: {}
};
export default LegalPersonReview;
