import { Stack, StackProps, useColorModeValue } from '@chakra-ui/react';

interface FormLayoutProps extends StackProps {}
const FormLayout: React.FC<FormLayoutProps> = (props) => {
  return (
    <Stack
      spacing={3.5}
      align="start"
      bg={useColorModeValue('white', '#171923')}
      border="2px solid #E5EDF1"
      borderRadius="10px"
      padding={{ base: 3, md: 9 }}
      {...props}
    />
  );
};

export default FormLayout;
