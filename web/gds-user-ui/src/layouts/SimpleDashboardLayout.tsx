import { Box, Stack, Flex } from '@chakra-ui/react';

type SimpleDashboardLayout = {
  children: React.ReactNode;
};
export const SimpleDashboardLayout: React.FC<SimpleDashboardLayout> = ({ children }) => {
  return (
    <Flex
      direction="column"
      align="center"
      maxW={'100%'}
      bg={'#F7F8FC'}
      px={58}
      py={10}
      fontFamily={'Open Sans'}
      position={'relative'}
      minHeight={'100vh'}>
      <Box>{children}</Box>
    </Flex>
  );
};
