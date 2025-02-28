import React from "react";
import { render, screen, fireEvent, waitFor } from "@testing-library/react";
import "@testing-library/jest-dom";
import App from "../App";

beforeEach(() => {
    jest.restoreAllMocks(); // Reset mocks before each test
});

test("displays an error message when API request fails", async () => {
    render(<App />);

    fireEvent.change(screen.getByLabelText(/Number of Items/i), { target: { value: "100" } });

    jest.spyOn(global, "fetch").mockResolvedValueOnce(
        Promise.resolve(
            new Response(JSON.stringify({ message: "Failed to fetch results" }), {
                status: 400,
                headers: { "Content-Type": "application/json" },
            })
        )
    );

    fireEvent.click(screen.getByRole("button", { name: /Calculate/i }));

    await waitFor(() => {
        expect(screen.getByText("Failed to fetch results")).toBeInTheDocument();
    });
});