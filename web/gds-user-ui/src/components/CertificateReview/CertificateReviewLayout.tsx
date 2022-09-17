import { Box, Stack } from '@chakra-ui/react';
import { ReactNode } from 'react';

type CertificateReviewLayoutProps = {
  children: ReactNode;
};

function CertificateReviewLayout({ children }: CertificateReviewLayoutProps) {
  return (
    <Box
      border="1px solid #DFE0EB"
      fontFamily={'Open Sans'}
      color={'#252733'}
      fontSize={18}
      p={5}
      px={5}>
      <Stack>{children}</Stack>
    </Box>
  );
}

export default CertificateReviewLayout;
