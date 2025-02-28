// src/App.tsx
import React, { useState, useEffect } from "react";
import { Container, Typography, TextField, Button } from "@mui/material";
import PackSizeItem from "./components/PackSizeItem";
import PackSizeForm from "./components/PackSizeForm";
import ErrorSnackbar from "./components/ErrorSnackbar";
import CalculationResultDialog from "./components/CalculationResultDialog";
import { fetchPackSizes, calculatePacks } from "./api";
import "./styles/styles.css";

const App = () => {
    const [items, setItems] = useState("");
    const [error, setError] = useState("");
    const [packSizes, setPackSizes] = useState<number[]>([]);
    const [result, setResult] = useState<{ pack: number; count: number }[]>([]);
    const [itemsOrdered, setItemsOrdered] = useState<number | null>(null);
    const [totalItemsUsed, setTotalItemsUsed] = useState<number | null>(null);
    const [openDialog, setOpenDialog] = useState(false);

    useEffect(() => {
        const loadPackSizes = async () => {
            try {
                const data = await fetchPackSizes();
                setPackSizes(data.pack_sizes || []);
            } catch (err) {
                setError((err as Error).message);
            }
        };
        loadPackSizes();
    }, []);

    const handleCalculate = async () => {
        if (!items || isNaN(Number(items))) {
            setError("Please enter a valid number of items.");
            return;
        }
        try {
            const data = await calculatePacks(Number(items));
            setResult(data.result);
            setItemsOrdered(data.itemsOrdered);
            setTotalItemsUsed(data.totalItemsUsed);
            setOpenDialog(true);
        } catch (err) {
            setError((err as Error).message);
        }
    };

    return (
        <Container maxWidth="sm" className="app-container">
            <Typography variant="h4" className="header">
                Pack Calculator
            </Typography>

            <TextField
                fullWidth
                label="Number of Items"
                variant="outlined"
                value={items}
                onChange={(e) => setItems(e.target.value)}
                className="input-field"
            />

            <Button
                variant="contained"
                color="primary"
                onClick={handleCalculate}
                className="calculate-btn"
                sx={{ mt: 2 }}
            >
                Calculate
            </Button>

            <CalculationResultDialog
                open={openDialog}
                onClose={() => setOpenDialog(false)}
                result={result}
                itemsOrdered={itemsOrdered}
                totalItemsUsed={totalItemsUsed}
            />

            {/* New Pack Sizes Grid */}
            <div className="pack-size-grid">
                {packSizes.map((size) => (
                    <PackSizeItem
                        key={size}
                        size={size}
                        onDelete={(deletedSize) =>
                            setPackSizes((prev) =>
                                prev.filter((s) => s !== deletedSize)
                            )
                        }
                    />
                ))}
            </div>

            <PackSizeForm onPackSizeUpdate={(newPackSizes) => setPackSizes(newPackSizes)} />

            <ErrorSnackbar error={error} setError={setError} />
        </Container>
    );
};

export default App;