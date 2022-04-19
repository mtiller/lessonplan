{% macro commentary(now) %}

Our solution for $x$ was {{ now.x }} and our solution for $y$ depends on time. Currently, $t={{ now.time.toFixed(2) }}$ which means that $y={{ now.y.toFixed(2) }}$.

{% endmacro %}
