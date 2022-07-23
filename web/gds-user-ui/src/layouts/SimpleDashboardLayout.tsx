import { Box, Stack, Flex } from '@chakra-ui/react';
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
          maxW={'100%'}
          bg={'#F7F8FC'}
          px={58}
          py={10}
          fontFamily={'Open Sans'}
          position={'relative'}
          minHeight={'100vh'}>
          <Box>{children}</Box>
        </Flex>
      ) : (
        children
      )}
    </>
  );
};
