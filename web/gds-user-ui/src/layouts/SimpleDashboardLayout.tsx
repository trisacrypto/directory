import { Box, Container, Flex, useColorModeValue } from '@chakra-ui/react';
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
          m="auto"
          bg={'#F7F8FC'}
          w="100%"
          px={58}
          py={10}
          fontFamily={'Open Sans'}
          position={'relative'}
          minHeight={'100vh'}>
          <Container maxW="7xl">{children}</Container>
        </Flex>
      ) : (
        <DashboardLayout>{children}</DashboardLayout>
      )}
    </>
  );
};
