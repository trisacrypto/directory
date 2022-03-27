import { Box, Stack } from '@chakra-ui/react';

type SimpleDashboardLayout = {
  children: React.ReactNode;
};
export const SimpleDashboardLayout: React.FC<SimpleDashboardLayout> = ({ children }) => {
  return (
    <Stack px={58}>
      <Box>{children}</Box>
    </Stack>
  );
};
