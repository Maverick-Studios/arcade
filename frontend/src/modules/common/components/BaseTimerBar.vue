<template>
    <div class="base-timer">
        <svg
            class="base-timer__svg"
            viewBox="0 0 50 5"
            xmlns="http://www.w3.org/2000/svg"
        >
            <path class="base-timer__path-elapsed" d="M 1, 1 L 49, 1"></path>
            <path
                :stroke-dasharray="dasharray"
                class="base-timer__path-remaining"
                :class="remainingPathColor"
                d="M 1, 1 L 49, 1"
            ></path>
        </svg>
        <span class="base-timer__label">
            {{ timeLeft }}
        </span>
    </div>
</template>

<script>
import { EventBus } from "@/eventBus";
import { Event } from "@/events";

const FULL_DASH_ARRAY = 48;
const WARNING_THRESHOLD = 5;
const DANGER_THRESHOLD = 2;
const COLOR_CODES = {
    healthy: {
        color: "green",
    },
    warning: {
        color: "orange",
        threshold: WARNING_THRESHOLD,
    },
    danger: {
        color: "red",
        threshold: DANGER_THRESHOLD,
    },
};
export default {
    name: "BaseTimerBar",
    props: {
        size: String,
    },

    data() {
        return {
            timeLimit: 0,
            timePassed: 0,
            timerInterval: null,
        };
    },
    computed: {
        dasharray() {
            return `${(this.timeFraction * FULL_DASH_ARRAY).toFixed(0)} 48`;
        },
        timeFraction() {
            return this.timeLeft / this.timeLimit;
        },
        remainingPathColor() {
            const { healthy, warning, danger } = COLOR_CODES;
            if (this.timeLeft <= DANGER_THRESHOLD) {
                return danger.color;
            } else if (this.timeLeft <= WARNING_THRESHOLD) {
                return warning.color;
            } else {
                return healthy.color;
            }
        },
        timeLeft() {
            return this.timeLimit - this.timePassed;
        },
    },
    watch: {
        timeLeft: function(newValue) {
            if (newValue === 0) {
                this.stopTimer();
            }
        },
    },
    methods: {
        startTimer() {
            this.timerInterval = setInterval(
                () => (this.timePassed += 1),
                1000
            );
        },
        stopTimer() {
            clearInterval(this.timeInterval)
        },
        resetTimer(data) {
            this.timeLimit = data;
            this.timePassed = 0;
        }
    },
    created() {
         EventBus.$on(Event.TIMER_RESET, (data) => {
            this.resetTimer(data);
        });

        this.timeLimit = this.$scribbleStoreService.getTimerDuration();
        this.startTimer();
    },
};
</script>

<style scoped lang="scss">
.base-timer {
    &__path-elapsed {
        stroke-width: 1px;
        stroke-linecap: round;
        stroke: grey;
    }
    &__path-remaining {
        stroke-width: 1px;
        stroke-linecap: round;
        transform-origin: right;
        transition: 1s linear all;
        stroke: currentColor;
        fill-rule: nonzero;
        &.green {
            color: rgb(65, 184, 131);
        }
        &.orange {
            color: orange;
        }
        &.red {
            color: red;
        }
    }
}
</style>
