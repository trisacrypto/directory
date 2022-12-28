import React from 'react';
import { Text } from '@chakra-ui/react';
import { getFormattedAmount } from 'utils/utils';
import type { ITrixo } from 'modules/dashboard/certificate/entities';

interface ComplianceThresholdProps {
  data:
    | Pick<
        ITrixo,
        'compliance_threshold' | 'compliance_threshold_currency' | 'must_comply_travel_rule'
      >
    | undefined;
}
const ComplianceThresholdRow = (data: ComplianceThresholdProps) => {
  const { compliance_threshold, compliance_threshold_currency, must_comply_travel_rule } =
    data?.data as any;
  const shouldShowComplianceThreshold = compliance_threshold || compliance_threshold !== 0;

  const getAmount = () => {
    return getFormattedAmount(compliance_threshold, compliance_threshold_currency);
  };
  const mustComplyTravelRule = () => {
    if (must_comply_travel_rule && compliance_threshold === 0) {
      return `${getFormattedAmount(compliance_threshold, compliance_threshold_currency)}`;
    }
    return 'N/A';
  };

  return (
    <>
      {shouldShowComplianceThreshold ? (
        <Text>{getAmount()}</Text>
      ) : (
        <Text>{mustComplyTravelRule()}</Text>
      )}
    </>
  );
};

export default ComplianceThresholdRow;
