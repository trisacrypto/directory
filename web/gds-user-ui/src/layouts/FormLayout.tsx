import { Stack } from "@chakra-ui/react";

type FormLayoutProps = {
  children: React.ReactNode;
};
const FormLayout: React.FC<FormLayoutProps> = ({ children }) => {
  return (
    <Stack
      spacing={3.5}
      align="start"
      border="2px solid #E5EDF1"
      borderRadius={2.5}
      padding={{ base: 3, md: 9 }}
    >
      {children}
    </Stack>
  );
};

export default FormLayout;
