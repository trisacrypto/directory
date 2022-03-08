import { Button, ButtonProps } from "@chakra-ui/react";

interface _ButtonProps extends ButtonProps {}

const FormButton: React.FC<_ButtonProps> = (props) => (
  <Button
    borderRadius={0}
    background="#55ACD8"
    color="#fff"
    _hover={{ background: "blue.200" }}
    {...props}
  />
);

export default FormButton;
