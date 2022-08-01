import { Box, Heading, Button } from '@chakra-ui/react';
import { Trans } from '@lingui/react';
import useCertificateStepper from 'hooks/useCertificateStepper';

type CertificateReviewHeaderProps = {
  step: number;
  title: string;
};

const CertificateReviewHeader = ({ step, title }: CertificateReviewHeaderProps) => {
  const { jumpToStep } = useCertificateStepper();

  return (
    <Box display={'flex'} justifyContent="space-between" pt={4} ml={0}>
      <Heading fontSize={20} mb="2rem">
        <Trans id={title}>{title}</Trans>
      </Heading>
      <Button
        bg={'blue'}
        color={'white'}
        height={'34px'}
        onClick={() => jumpToStep(step)}
        _hover={{
          bg: '#10aaed'
        }}>
        <Trans id="Edit">Edit</Trans>
      </Button>
    </Box>
  );
};

export default CertificateReviewHeader;
