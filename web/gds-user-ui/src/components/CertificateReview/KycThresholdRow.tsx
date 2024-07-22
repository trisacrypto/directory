import React, { Suspense } from 'react';
import { Text } from '@chakra-ui/react';
import { getFormattedAmount } from 'utils/utils';
import type { ITrixo } from 'modules/dashboard/certificate/entities';

interface KycThresholdRowProps {
  data:
    | Pick<ITrixo, 'kyc_threshold' | 'kyc_threshold_currency' | 'has_required_regulatory_program'>
    | undefined;
}
const KycThresholdRow = (data: KycThresholdRowProps) => {
  const { kyc_threshold, kyc_threshold_currency, has_required_regulatory_program } =
    data?.data as any;
  const shouldShowKycThreshold = kyc_threshold || +kyc_threshold !== 0;

  const getAmount = () => {
    return getFormattedAmount(kyc_threshold, kyc_threshold_currency);
  };
  const hasRequiredRegulatoryProgram = () => {
    if (has_required_regulatory_program === 'yes' && kyc_threshold === 0) {
      return `${getFormattedAmount(kyc_threshold, kyc_threshold_currency)}`;
    }
    return 'N/A';
  };

  return (
    <>
      <Suspense fallback={<div>Loading...</div>}>
        {shouldShowKycThreshold ? (
          <Text>{getAmount()}</Text>
        ) : (
          <Text>{hasRequiredRegulatoryProgram()}</Text>
        )}
      </Suspense>
    </>
  );
};

export default KycThresholdRow;
