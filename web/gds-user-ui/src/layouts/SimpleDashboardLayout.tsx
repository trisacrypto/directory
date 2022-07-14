import { Box, Flex, useColorModeValue } from '@chakra-ui/react';
import useAuth from 'hooks/useAuth';
type SimpleDashboardLayout = {
  children: React.ReactNode;
};
import DashboardLayout from './DashboardLayout';
export const SimpleDashboardLayout: React.FC<SimpleDashboardLayout> = ({ children }) => {
  const { isUserAuthenticated } = useAuth();
  const bg = useColorModeValue('#F7F8FC', 'gray.800');

  return (
    <>
      {!isUserAuthenticated ? (
        <Flex
          direction="column"
          align="center"
          maxW={'100%'}
          bg={bg}
          px={58}
          py={10}
          fontFamily={'Open Sans'}
          position={'relative'}
          minHeight={'100vh'}>
          <Box>{children}</Box>
        </Flex>
      ) : (
        <DashboardLayout>{children}</DashboardLayout>
      )}
    </>
  );
};
