import { Meta, Story } from "@storybook/react";
import DeleteButton from "./DeleteButton";

type DeleteButtonProps = {
  children: React.ReactNode;
};

export default {
  title: "components/DeleteButton",
  component: DeleteButton,
} as Meta;

const Template: Story<DeleteButtonProps> = (args) => <DeleteButton {...args} />;

export const Default = Template.bind({});
Default.args = {};
