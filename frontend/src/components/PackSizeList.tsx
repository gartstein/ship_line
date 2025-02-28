// src/components/PackSizeList.tsx
import React from "react";
import { Typography, Grid } from "@mui/material";

const PackSizeList = ({ packSizes }: { packSizes: number[] }) => {
    return (
        <>
            <Typography variant="subtitle1" className="title">
                Current Pack Sizes:
            </Typography>
            <Grid container spacing={1} className="pack-grid">
                {packSizes.length > 0 ? (
                    packSizes.map((size, index) => (
                        <Grid item xs={4} sm={3} md={2} key={index} className="pack-item">
                            <Typography className="pack-text">{size} items</Typography>
                        </Grid>
                    ))
                ) : (
                    <Typography className="no-pack">No pack sizes available.</Typography>
                )}
            </Grid>
        </>
    );
};

export default PackSizeList;