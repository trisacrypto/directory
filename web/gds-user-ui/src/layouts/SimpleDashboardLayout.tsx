import { Box, Container, Flex, useColorModeValue } from '@chakra-ui/react';
import useAuth from 'hooks/useAuth';
import { userSelector } from 'modules/auth/login/user.slice';
import { useSelector } from 'react-redux';
type SimpleDashboardLayout = {
  children: React.ReactNode;
};
export const SimpleDashboardLayout: React.FC<SimpleDashboardLayout> = ({ children }) => {
  const { isLoggedIn } = useSelector(userSelector);
  return (
    <>
      {!isLoggedIn ? (
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
        children
      )}
    </>
  );
};
