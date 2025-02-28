import "@testing-library/jest-dom";
import React from "react";
import { render, fireEvent } from "@testing-library/react";
import Button from "../Button";

test("renders button correctly", () => {
    const { getByText } = render(<Button onClick={() => {}} loading={false} />);
    expect(getByText("Calculate")).toBeInTheDocument();
});

test("calls onClick handler when clicked", () => {
    const handleClick = jest.fn();
    const { getByText } = render(<Button onClick={handleClick} loading={false} />);
    const buttonElement = getByText("Calculate");

    fireEvent.click(buttonElement);
    expect(handleClick).toHaveBeenCalledTimes(1);
});

test("shows 'Calculating...' when loading", () => {
    const { getByText } = render(<Button onClick={() => {}} loading={true} />);
    expect(getByText("Calculating...")).toBeInTheDocument();
});